import csv

with open('./ultimateActorNames.csv') as inp, open('./finalasitoActorasitoNamosito.csv', 'w') as output:
    writer = csv.writer(output)
    for line in inp:
        out = str(line).replace(',', '|').replace('|', ',', 1)
        print(out)
        writer.writerow([out])