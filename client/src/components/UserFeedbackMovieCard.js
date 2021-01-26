import { Card, CardWrapper } from "react-swipeable-cards";

export default function UserFeedbackMovieCard({ message }) {
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
      </Card>
    </CardWrapper>
  );
}

let stylingTitle = {
  textAlign: "center",
  fontWeight: "bold",
  fontSize: "35px",
  marginTop: "270px",
};
