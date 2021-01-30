import { LinearProgress } from "@material-ui/core";
import { Card, CardWrapper } from "react-swipeable-cards";

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
        <div style={stylingTitle}>{message}</div>
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
