import { LinearProgress, makeStyles, Typography } from "@material-ui/core";
import { Card, CardWrapper } from "@xhuliodo/react-swipeable-cards";

export default function DataStateMovieCard({ message, type }) {
  const classes = useStyle();
  
  return (
    <CardWrapper style={{ paddingTop: "0px" }}>
      <Card
        style={{
          backgroundSize: "contain",
          backgroundRepeat: "no-repeat",
          backgroundPosition: "center",
        }}
      >
        <Typography className={classes.title}>{message}</Typography>
        {type === "loading" ? (
          <LinearProgress className={classes.loading} />
        ) : null}
      </Card>
    </CardWrapper>
  );
}

const useStyle = makeStyles(() => ({
  title: {
    color: "black",
    textAlign: "center",
    fontWeight: "bold",
    fontSize: "35px",
    marginTop: "45%",
  },
  loading: {
    width: "85%",
    margin: "50px auto",
  },
}));

let stylingTitle = {};
