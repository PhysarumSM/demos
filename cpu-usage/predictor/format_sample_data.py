import csv
import sys

def format_data(filename, output_filename, num_columns):
    rows = [[]]
    with open(filename, 'r') as f:
        for line in f:
            value = float(line.strip())
            if len(rows[-1]) >= num_columns:
                rows.append([])
            rows[-1].append(value)

    if len(rows[-1]) < num_columns:
        del rows[-1]

    with open(output_filename, 'w') as f:
        cw = csv.writer(f)
        for row in rows:
            cw.writerow(row)

if __name__ == '__main__':
    if len(sys.argv) != 4:
        print('Usage: {} <input file> <output file> <num columns>'.format(sys.argv[0]))
        sys.exit(1)

    format_data(sys.argv[1], sys.argv[2], int(sys.argv[3]))