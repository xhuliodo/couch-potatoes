import csv

# List all titles and all corresponding actors
with open('principals.tsv') as inp, open('allActorsAllMovies.csv','w') as output:
    writer = csv.writer(output)
    for line in inp:
        splitRow = line.split()
        if(splitRow[3] == 'actor' or splitRow[3] == 'actress'):
            writer.writerow([splitRow[0], splitRow[2]])

# group all actors of each movie
oneMovieArray = []
with open('allActorsAllMovies.csv') as inp, open('actorsIDperMovie.csv', 'w') as output:
    writer = csv.writer(output)
    for line in inp:
        data = line.split(',')
        data[1] = data[1][:-2]
        if oneMovieArray:
            if oneMovieArray[0] == data[0]:
                oneMovieArray.append(data[1])
            else:
                writer.writerow(oneMovieArray)
                del oneMovieArray[:]
                oneMovieArray.append(data[0])
                oneMovieArray.append(data[1])
        else:
            oneMovieArray.append(data[0])
            oneMovieArray.append(data[1])
            

# create a list of titleID -> actor names

IDtoNameDict = {}
with open('./name.tsv') as f:
    for line in f:
        IDtoNameDict[line.split('\t')[0]] = line.split('\t')[1]
tempActors = []
with open('./actorsIDperMovie.csv') as inp, open('./actorsNamePerMovie.csv', 'w') as output:
    writer = csv.writer(output, delimiter=',')
    for line in inp:
        data = line.split(',')
        tempActors.append(data[0])
        for actor in data[1:]:
            try:
                tempActors.append(IDtoNameDict[actor])
            except:
                print('not-found')
        writer.writerow(tempActors)
        del tempActors[:]

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
# with open('./actorsNamePerMovie.csv') as inp, open ('./ultimateActorNames.csv', 'w') as output:
#     writer = csv.writer(output)
#     for line in inp:
#         if (line.split(',')[0] in usedTitleIDs and len(line.split(',')) > 1):
#             writer.writerow([line])