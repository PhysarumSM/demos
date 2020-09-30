import flask
import threading
import sys
import time

app = flask.Flask(__name__)

data = []

@app.route('/upload', methods=['POST'])
def upload():
    print('Upload post request from:', flask.request.remote_addr)
    posted_data = flask.request.json
    print(posted_data)
    data.extend(posted_data)
    return 'OK'

# Train model
# For now, use dummy model that takes average of inputs
def run_training():
    print('Training on data:', data)
    total_error = 0.0
    for sample in data:
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
    if len(sys.argv) < 2:
        print('Usage: {} <server port>'.format(sys.argv[0]))
        sys.exit()

    port = int(sys.argv[1])
    threading.Thread(target=app.run, kwargs={'host': '0.0.0.0', 'port': port},
            daemon=True).start()

    while True:
        run_training()
        time.sleep(1)