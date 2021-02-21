import { useAuth0 } from "@auth0/auth0-react";
import { Button, Snackbar } from "@material-ui/core";
import { useEffect, useState } from "react";
import { useHistory } from "react-router-dom";

export default function WelcomePage() {
  const history = useHistory();
  const redirect = () => {
    history.push("/solo");
  };

  const action = (
    <Button variant="contained" onClick={redirect} size="small">
      GO
    </Button>
  );

  const { isAuthenticated } = useAuth0();

  console.log(isAuthenticated, typeof isAuthenticated);

  const [open, setOpen] = useState(isAuthenticated);

  useEffect(() => {
    setOpen(isAuthenticated);
  }, [isAuthenticated]);

  const handleClose = (event, reason) => {
    if (reason === "clickaway") {
      return;
    }

    setOpen(false);
  };

  return (
    <div>
      <Snackbar
        open={open}
        onClose={handleClose}
        action={action}
        message="Go back to using the app?"
        // anchorOrigin={{ horizontal: "center", vertical: "top" }}
      />
      This is the page where i gotta sell the product to new users so i can make
      them log in
    </div>
  );
}
