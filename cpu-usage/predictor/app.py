import flask
import multiprocessing
import multiprocessing.managers
import sys
import threading
import time

app = flask.Flask(__name__)

data = []
lock = multiprocessing.Lock()
manager_port = 4300 # arbitrary number which will be set later

def put_data(new_data):
    with lock:
        data.extend(new_data)

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

# Train model
# For now, use dummy model that takes average of inputs
def run_training():
    manager = multiprocessing.managers.BaseManager(('127.0.0.1', manager_port), b'password')
    manager.register('get_data')
    manager.connect()
    current_data = manager.get_data()._getvalue()
    print('Training on data:', current_data)
    total_error = 0.0
    for sample in current_data:
        inputs = sample[:4]
        target = sample[4]
        prediction = 0.0
        for x in inputs:
            prediction += x
        prediction = prediction / 4.0
        total_error += abs(target - prediction)
        # print('inputs:', inputs, ', target:', target, ', prediction:', prediction)
    print('Total Error:', total_error)

if __name__ == '__main__':
    if len(sys.argv) < 3:
        print('Usage: {} <proxy port> <server port>'.format(sys.argv[0]))
        sys.exit()

    proxy_port = int(sys.argv[1])
    server_port = int(sys.argv[2])

    manager_port = max(proxy_port, server_port) + 1
    manager = multiprocessing.managers.BaseManager(('127.0.0.1', manager_port), b'password')
    manager.register('put_data', put_data)
    manager.register('get_data', get_data)
    manager.start()

    threading.Thread(target=app.run, kwargs={'host': '0.0.0.0', 'port': server_port},
            daemon=True).start()

    while True:
        run_training()
        time.sleep(1)

    manager.end()