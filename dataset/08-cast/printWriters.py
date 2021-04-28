import csv

# giving credits where it's due

# # create a list of titleID -> director names

# IDtoNameDict = {}

# with open('./name.tsv') as f:
#     for line in f:
#         IDtoNameDict[line.split('\t')[0]] = line.split('\t')[1]


# tempWriters = []
# with open('./crew.tsv') as inp, open('./writerNames.csv', 'w') as output:
#     writer = csv.writer(output, delimiter=',')
#     for line in inp:
#         tempWriters.append(line.split()[0])
#         writers = line.split()[2].split(',')
#         if(writers == '\N'):
#             writer.writerow('not-found')
#         else:
#             for oneWriter in writers:
#                 try:
#                     tempWriters.append(IDtoNameDict[oneWriter])
#                 except:
#                     tempWriters.append('not-found')
#             writer.writerow(tempWriters)
#             del tempWriters[:]

# create a list of titleIDs in actual dataset
usedTitleIDs = []
with open('./movies.csv') as inp:
    for line in inp:
        try:
            startIndex = line.index('/tt')
            usedTitleIDs.append(line[startIndex + 1 : startIndex + 10])
        except:
            print('err')       
inp.close()

tempArr = []
with open('./writerNames.csv') as inp, open ('./finalWriterNames.csv', 'w') as output:
    writer = csv.writer(output, delimiter=',')
    for line in inp:
        for item in line.split(','):
            tempArr.append(item)
        if(tempArr[0] in usedTitleIDs):
            writer.writerow(tempArr)
        del tempArr[:]
        



