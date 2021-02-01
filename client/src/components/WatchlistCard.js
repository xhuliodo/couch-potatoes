import {
  makeStyles,
  Card,
  CardContent,
  CardMedia,
  IconButton,
  Typography,
  CardActions,
  Link,
  Tooltip,
} from "@material-ui/core";
import { Link as LinkIcon, ThumbDown, ThumbUp } from "@material-ui/icons";
import { useState } from "react";

import { useMutation } from "react-query";
import { rateMovie } from "../utils/requests";

export default function MovieCardWatchlist({ m }) {
  const classes = useStyles();

  const [movies, setMovies] = useState(m);

  const rate = useMutation((mutationData) => rateMovie(mutationData));

  const handleRate = (movieId, action) => {
    const mutationData = {
      movieId,
      userId: 2,
      action,
      successFunc: () => console.log("eeeee"),
    };
    rate.mutate(mutationData);
    const unratedMovies = movies.filter((m) => m.movieId !== movieId);
    setMovies(unratedMovies);
  };

  // const successFunc = (movieId) => {

  // };

  return movies.map((m) => {
    const image = `${m.posterUrl}`;
    return (
      <Card key={m.movieId} className={classes.root}>
        <CardMedia className={classes.cover} image={image} />
        <div className={classes.details}>
          <CardContent className={classes.content}>
            <Typography component="h5" variant="h5">
              {m.title}
            </Typography>
            <Typography variant="subtitle1" color="textSecondary">
              {m.releaseYear}
            </Typography>
          </CardContent>
          <CardActions disableSpacing>
            <IconButton
              onClick={() => {
                handleRate(m.movieId, "love");
              }}
              variant="contained"
              color="primary"
            >
              <ThumbUp className={classes.rateIcons} fontSize="inherit" />
            </IconButton>
            <IconButton
              onClick={() => handleRate(m.movieId, "hate")}
              variant="contained"
              color="secondary"
            >
              <ThumbDown className={classes.rateIcons} fontSize="inherit" />
            </IconButton>
            <div className={classes.delete}>
              <Tooltip title="IMDB Link" arrow placement="top">
                <IconButton variant="contained">
                  <Link
                    href={m.imdbLink}
                    rel="noreferrer"
                    style={{ height: "24px" }}
                    target="_blank"
                  >
                    <LinkIcon fontSize="inherit" />
                  </Link>
                </IconButton>
              </Tooltip>
            </div>
          </CardActions>
        </div>
      </Card>
    );
  });
}

const useStyles = makeStyles((theme) => ({
  root: {
    display: "flex",
    margin: "20px 0",
  },
  details: {
    display: "flex",
    flexDirection: "column",
    width: "90%",
  },
  cover: {
    width: 125,
  },
  rateIcons: {
    height: 32,
    width: 32,
  },
  delete: {
    width: "100%",
    display: "flex",
    justifyContent: "flex-end",
  },
}));
