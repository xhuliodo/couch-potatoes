import { LinearProgress, Typography } from "@material-ui/core";
import { Card, CardWrapper } from "@xhuliodo/react-swipeable-cards";

export default function UserFeedbackMovieCard({ message, type }) {
  return (
    <CardWrapper style={{ paddingTop: "0px" }}>
      <Card
        style={{
          backgroundSize: "contain",
          backgroundRepeat: "no-repeat",
          backgroundPosition: "center",
        }}
      >
        <Typography style={stylingTitle}>{message}</Typography>
        {type === "loading" ? (
          <LinearProgress style={{ width: "85%", margin: "50px auto" }} />
        ) : null}
      </Card>
    </CardWrapper>
  );
}

let stylingTitle = {
  textAlign: "center",
  fontWeight: "bold",
  fontSize: "35px",
  marginTop: "45%",
};
