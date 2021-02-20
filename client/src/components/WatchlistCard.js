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
  SentimentDissatisfied,
  SentimentVerySatisfiedRounded,
} from "@material-ui/icons";

export default function MovieCardWatchlist({
  m,
  handleRate,
  handleRemove,
  deleted,
  animation,
}) {
  const classes = useStyles();
  const theme = useTheme();

  const image = `http://thumb.cp.dev.cloudapp.al/thumbnail_${m.movieId}.jpg`;

  const rating = m?.rating;

  return (
    <Slide
      // appear
      // enter={false}
      timeout={{ enter: 0, exit: 300 }}
      
      // enter={false}
      direction={animation}
      in={!deleted?.includes(m.movieId)}
      unmountOnExit
    >
      <Card elevation={5} className={classes.root}>
        <CardMedia className={classes.cover} image={image} />
        <div className={classes.details}>
          <CardContent className={classes.content}>
            <Typography style={{ display: "flex" }} component="h5" variant="h5">
              {m.title}
              {/* <div className={classes.link}> */}
              <Tooltip title="IMDB Link" arrow placement="top">
                <Link
                  href={m.imdbLink}
                  rel="noreferrer"
                  style={{ marginLeft: "10px" }}
                  target="_blank"
                >
                  <Info className={classes.infoIcon} />
                </Link>
              </Tooltip>
              {/* </div> */}
            </Typography>
            <Typography variant="subtitle1" color="textSecondary">
              {m.releaseYear}
            </Typography>
          </CardContent>
          <CardActions disableSpacing>
            <Button
              disabled={rating === undefined ? false : true}
              style={
                rating === 0
                  ? {
                      backgroundColor: theme.palette.secondary.main,
                      margin: "0 10px",
                    }
                  : { margin: "0 10px" }
              }
              onClick={() => handleRate(m.movieId, "hate")}
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
                      margin: "0 10px",
                    }
                  : { margin: "0 10px" }
              }
              disabled={rating === undefined ? false : true}
              onClick={() => handleRate(m.movieId, "love")}
              variant="contained"
              color="primary"
            >
              <SentimentVerySatisfiedRounded
                className={classes.rateIcons}
                fontSize="inherit"
              />
            </Button>

            <div className={classes.trash}>
              <IconButton
                onClick={() => handleRemove(m.movieId)}
                variant="contained"
              >
                <Delete color="secondary" />
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
    // margin: "20px 0",
    marginBottom: "20px",
    marginTop: "5px",
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
    height: 26,
    width: 26,
  },
  link: {
    width: "100%",
    display: "flex",
    justifyContent: "flex-end",
  },
  trash: {
    width: "100%",
    display: "flex",
    justifyContent: "flex-end",
  },
  infoIcon: {
    height: 20,
    width: 20,
  },
}));
