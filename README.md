[![linkedin-profile](https://camo.githubusercontent.com/5b6f22d22d8701b7cf49ef2526af05b5a17695683a8ced9efda2439591cd0cdf/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f4c696e6b6564496e2d50726f66696c652d696e666f726d6174696f6e616c3f7374796c653d666c6174266c6f676f3d6c696e6b6564696e266c6f676f436f6c6f723d776869746526636f6c6f723d304437364138)](https://www.linkedin.com/in/xhulio-doda-745b41164/)
[![email](https://img.shields.io/badge/email-xhuliodo-red)]([mailto](mailto:xhuliodo@gmail.com)
# Couch Potatoes: Full-stack project for providing real time movie recommendations

Most people hear of new movies either from friends or from online articles reviewing them. After getting a wind of an interesting movie, you have to do the homework researching who's the director, writer or who's in the cast and then asking friends with similar movies, who might have watched it, if the movie it's worth giving a try. To simplify this process, Couch Potatoes was made to provide personalised movie recommendations.

The personalised part was made possible using a graph database (Neo4j), following two strategies:
- Content-based (similarity between movies)
- User-based (similarity between users)

The API was build using Golang, following the Clean Architecture principles laid by Uncle Bob.

Interfaces were built using React and the main component is inspired by tinder's swiping cards.

## Getting started

After cloning the repo, all you have to do is execute **docker-compose up**. 
After that you have to wait for a couple of minutes, till the seeding process finishes (the dataset used is the Small MovieLens Latest Dataset [2021]). 
(To get the full experience, you might also want to get the thumbnails which you can find the steps in the /dataset folder. It's a crawling script that downloads the thumbnails from imdb's site.)

That's all, enjoy swiping and getting some good movie recommendations.

contact: xhuliodo@gmail.com
