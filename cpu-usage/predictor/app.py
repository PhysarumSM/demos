import flask
import multiprocessing
import multiprocessing.managers
import socket
import sys
import threading
import time

import model

app = flask.Flask(__name__)

data = model.read_csv_data('sample_cpu_data.csv')
lock = multiprocessing.Lock()
manager_port = 4300 # arbitrary number which will be set later

def put_data(new_data):
    with lock:
        for row in new_data:
            if len(row) == 5:
                data.append(row)

def get_data():
    with lock:
        return data

@app.route('/upload', methods=['POST'])
def upload():
    print('Upload post request from:', flask.request.remote_addr)
    posted_data = flask.request.json
    print(posted_data)
    manager = multiprocessing.managers.BaseManager(('127.0.0.1', manager_port), b'password')
    manager.register('put_data')
    manager.connect()
    manager.put_data(posted_data)
    return 'OK'

@app.route('/data', methods=['GET'])
def get():
    print('Get data request from:', flask.request.remote_addr)
    manager = multiprocessing.managers.BaseManager(('127.0.0.1', manager_port), b'password')
    manager.register('get_data')
    manager.connect()
    current_data = manager.get_data()._getvalue()
    return flask.json.jsonify(current_data)

# Train model
# For now, use dummy model that takes average of inputs
def run_training():
    manager = multiprocessing.managers.BaseManager(('127.0.0.1', manager_port), b'password')
    manager.register('get_data')
    manager.connect()
    current_data = manager.get_data()._getvalue()
    model.train_variations(current_data)

def next_free_port(exclude_ports, port=1024, max_port=65535):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    while port <= max_port:
        if port in exclude_ports:
            continue
        try:
            sock.bind(('', port))
            sock.close()
            return port
        except OSError:
            port += 1
    raise IOError('no free ports')

if __name__ == '__main__':
    if len(sys.argv) < 3:
        print('Usage: {} <proxy port> <server port>'.format(sys.argv[0]))
        sys.exit()

    proxy_port = int(sys.argv[1])
    server_port = int(sys.argv[2])

    manager_port = next_free_port({proxy_port, server_port})
    manager = multiprocessing.managers.BaseManager(('127.0.0.1', manager_port), b'password')
    manager.register('put_data', put_data)
    manager.register('get_data', get_data)
    manager.start()
    print('Run multiprocessing manager on 127.0.0.1:{}'.format(manager_port))

    threading.Thread(target=app.run, kwargs={'host': '0.0.0.0', 'port': server_port},
            daemon=True).start()

    while True:
        run_training()
        time.sleep(1)

    manager.end()