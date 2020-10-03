import csv
import numpy
import torch
import torch.nn as nn

class CpuUsagePredictorSimple(nn.Module):
    def __init__(self, num_inputs=4):
        super().__init__()
        self.num_inputs = num_inputs
        self.name = 'CpuUsagePredictorSimple'
        self.fc1 = nn.Linear(num_inputs, 1)

    def forward(self, x):
        x = x.view(-1, self.num_inputs)
        x = self.fc1(x)
        x = x.squeeze(1) # Flatten to [batch_size]
        return x

class CpuUsagePredictor(nn.Module):
    def __init__(self, num_inputs=4):
        super().__init__()
        self.num_inputs = num_inputs
        self.name = 'CpuUsagePredictor'
        self.fc1 = nn.Linear(num_inputs, num_inputs)
        self.fc2 = nn.Linear(num_inputs, 1)

    def forward(self, x):
        x = x.view(-1, self.num_inputs)
        x = self.fc1(x)
        x = self.fc2(x)
        x = x.squeeze(1) # Flatten to [batch_size]
        return x

class CpuUsagePredictorComplex(nn.Module):
    def __init__(self, num_inputs=4):
        super().__init__()
        self.num_inputs = num_inputs
        self.name = 'CpuUsagePredictorComplex'
        self.fc1 = nn.Linear(num_inputs, num_inputs)
        self.fc2 = nn.Linear(num_inputs, num_inputs//2)
        self.fc3 = nn.Linear(num_inputs//2, 1)

    def forward(self, x):
        x = x.view(-1, self.num_inputs)
        x = self.fc1(x)
        x = self.fc2(x)
        x = self.fc3(x)
        x = x.squeeze(1) # Flatten to [batch_size]
        return x

def train(model, train_data, epochs=20, lr=0.01):
    criterion = torch.nn.MSELoss()
    optimizer = torch.optim.Adam(model.parameters(), lr=lr)
    losses = []
    for _ in range(epochs):
        epoch_losses = []
        for inputs, targets in zip(train_data[0], train_data[1]):
            optimizer.zero_grad()
            outputs = model(inputs)
            loss = criterion(outputs, targets)
            loss.backward()
            optimizer.step()
            epoch_losses.append(loss.item())
        avg_loss = numpy.mean(epoch_losses)
        losses.append(avg_loss)
    print('Final loss:', losses[-1], ', Min loss:', numpy.min(losses))
    return losses

def average_predictor(train_data):
    criterion = torch.nn.MSELoss()
    losses = []
    for inputs, targets in zip(train_data[0], train_data[1]):
        outputs = torch.mean(inputs, 0, keepdim=True)
        loss = criterion(outputs, targets)
        losses.append(loss.item())
    avg_loss = numpy.mean(losses)
    print('Final loss:', avg_loss)
    return losses

def data_list_to_tensors(data_list, num_inputs=4):
    data_tensor = torch.Tensor(data_list)
    inputs = data_tensor[:, :num_inputs]
    targets = data_tensor[:, num_inputs:num_inputs+1]
    return (inputs, targets)

def concat_squared_inputs(inputs):
    squared = inputs**2
    new_inputs = torch.cat((inputs, squared), dim=1)
    return new_inputs

def train_variations(data_list):
    inputs, targets = data_list_to_tensors(data_list)
    inputs2 = concat_squared_inputs(inputs)
    datasets = [
        (inputs, targets),
        (inputs2, targets)
    ]
    model_classes = [
        CpuUsagePredictorSimple,
        CpuUsagePredictor,
        CpuUsagePredictorComplex
    ]

    for dataset in datasets:
        print('Try average predictor on data of size', dataset[0].size())
        average_predictor(dataset)
        for ModelClass in model_classes:
            num_inputs = dataset[0].size()[1]
            model = ModelClass(num_inputs=num_inputs)
            print('Training {} on data of size {}'.format(model.name, dataset[0].size()))
            train(model, dataset)

def read_csv_data(filename):
    csv_data = []
    with open(filename, 'r') as f:
        cr = csv.reader(f)
        csv_data = [[float(value) for value in row] for row in cr]
    return csv_data

if __name__ == '__main__':
    sample_data = read_csv_data('sample_cpu_data.csv')
    train_variations(sample_data)
