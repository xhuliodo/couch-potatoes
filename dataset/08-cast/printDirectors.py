import csv

# create a list of titleID -> director names

IDtoNameDict = {}

with open('./names.tsv') as f:
    for line in f:
        IDtoNameDict[line.split('\t')[0]] = line.split('\t')[1]


tempDirectors = []
with open('./crew.tsv') as inp, open('./directorNames.csv', 'w') as output:
    writer = csv.writer(output, delimiter=',')
    for line in inp:
        tempDirectors.append(line.split()[0])
        directors = line.split()[1].split(',')
        if(directors == '\N'):
            writer.writerow('not-found')
        else:
            for director in directors:
                try:
                    tempDirectors.append(IDtoNameDict[director])
                except:
                    tempDirectors.append('not-found')
            writer.writerow(tempDirectors)
            del tempDirectors[:]

# # create a list of titleIDs in actual dataset
# usedTitleIDs = []
# with open('./movies.csv') as inp:
#     for line in inp:
#         try:
#             startIndex = line.index('/tt')
#             usedTitleIDs.append(line[startIndex + 1 : startIndex + 10])
#         except:
#             print('err')       
# inp.close()

# tempArr = []
# with open('./directorNames.csv') as inp, open ('./finalDirectorNames.csv', 'w') as output:
#     writer = csv.writer(output, delimiter=',')
#     for line in inp:
#         for item in line.split(','):
#             tempArr.append(item)
#         if(tempArr[0] in usedTitleIDs):
#             writer.writerow(tempArr)
#         del tempArr[:]
        



