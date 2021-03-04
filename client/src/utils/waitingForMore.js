import { LinearProgress, Typography } from "@material-ui/core";

export const waitingForMore = () => (
  <div>
    <Typography
      style={{
        color: "black",
        textAlign: "center",
        fontWeight: "bold",
        fontSize: "35px",
        marginTop: "45%",
      }}
    >
      Fetching movies...
    </Typography>
    <LinearProgress style={{ width: "85%", margin: "50px auto" }} />
  </div>
);
