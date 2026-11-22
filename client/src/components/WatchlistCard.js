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
  useTheme,
  Button,
  Slide,
} from "@material-ui/core";
import {
  Delete,
  Info,
  MovieOutlined,
  SentimentDissatisfied,
  SentimentVerySatisfiedRounded,
} from "@material-ui/icons";
import { thumbnailUrl } from "../utils/config";
import usePosterFallback from "./usePosterFallback";

export default function MovieCardWatchlist({
  m,
  handleRate,
  handleRemove,
  deleted,
  animation,
  lastElementRef,
}) {
  const classes = useStyles();
  const theme = useTheme();

  const image = m.movie.movieId ? thumbnailUrl(m.movie.movieId) : null;
  const { showFallback: posterFailed, onError: handlePosterError } =
    usePosterFallback(image);
  const rating = m?.rating;

  return (
    <Slide
      timeout={{ enter: 0, exit: 250 }}
      direction={animation}
      in={!deleted?.includes(m.movie.movieId)}
      unmountOnExit
    >
      <Card ref={lastElementRef} elevation={3} className={classes.root}>
        {posterFailed ? (
          <div
            className={`${classes.cover} ${classes.posterless}`}
            aria-label="Poster unavailable"
          >
            <MovieOutlined className={classes.posterlessIcon} />
          </div>
        ) : (
          <CardMedia className={classes.cover} image={image} />
        )}
        {!posterFailed && (
          <img
            src={image}
            alt=""
            onError={handlePosterError}
            style={{ display: "none" }}
          />
        )}
        <div className={classes.details}>
          <div style={{ display: "flex", flexDirection: "row" }}>
            <CardContent style={{ flex: "auto" }}>
              <Typography>
                <strong>{m.movie.title}</strong>
              </Typography>
              <Typography variant="subtitle1" color="textSecondary">
                {m.movie.releaseYear}
              </Typography>
            </CardContent>
            <div className={classes.infoDiv}>
              <Tooltip title="IMDB Link" arrow placement="top">
                <Link href={m.movie.moreInfoLink} rel="noreferrer" target="_blank">
                  <Info color="action" className={classes.infoIcon} />
                </Link>
              </Tooltip>
            </div>
          </div>
          <CardActions>
            <Button
              disabled={rating === undefined ? false : true}
              className={classes.disabledButton}
              style={
                rating === 0
                  ? {
                      backgroundColor: theme.palette.secondary.main,
                    }
                  : null
              }
              onClick={() => handleRate(m.movie.movieId, "hate")}
              variant="contained"
              color="secondary"
            >
              <SentimentDissatisfied
                className={classes.rateIcons}
                fontSize="inherit"
              />
            </Button>
            <Button
              style={
                rating === 1
                  ? {
                      backgroundColor: theme.palette.primary.main,
                    }
                  : null
              }
              disabled={rating === undefined ? false : true}
              onClick={() => handleRate(m.movie.movieId, "love")}
              variant="contained"
              color="primary"
              className={classes.buttonContainer}
            >
              <SentimentVerySatisfiedRounded
                className={classes.rateIcons}
                fontSize="inherit"
              />
            </Button>

            <div className={classes.trash}>
              <IconButton
                onClick={() => handleRemove(m.movie.movieId)}
                variant="contained"
              >
                <Delete className={classes.trashColor} />
              </IconButton>
            </div>
          </CardActions>
        </div>
      </Card>
    </Slide>
  );
}

const useStyles = makeStyles((theme) => ({
  root: {
    display: "flex",
    marginBottom: "20px",
    marginTop: "5px",
  },
  details: {
    display: "flex",
    flexDirection: "column",
    width: "90%",
  },
  cover: {
    minWidth: 50,
    width: 110,
  },
  posterless: {
    background: "linear-gradient(135deg, #505874 0%, #8e72a7 100%)",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    color: "#fff",
  },
  posterlessIcon: {
    fontSize: 32,
    opacity: 0.8,
  },
  rateIcons: {
    height: 24,
    width: 24,
    "@media (max-width:330px)": {
      height: 16,
      width: 16,
    },
  },
  trash: {
    width: "100%",
    display: "flex",
    justifyContent: "flex-end",
  },
  trashColor: {
    color: "#a4a4a4",
  },
  infoDiv: {
    display: "flex",
    justifyContent: "center",
    marginTop: "18px",
    marginRight: "22px",
    height: "fit-content",
  },
  infoIcon: {
    height: 20,
    width: 20,
    // color: "#505874",
    // color: "#a4a4a4",
  },
  disabledButton: {
    marginLeft: 8,
    minWidth: 36,
    width: "15vw",
    maxWidth: 70,
  },
  buttonContainer: {
    minWidth: 36,
    width: "15vw",
    maxWidth: 70,
  },
}));
