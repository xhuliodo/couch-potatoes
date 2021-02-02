import { useEffect, useRef, useState } from "react";

// material ui components
import { makeStyles } from "@material-ui/core/styles";
import {
  Avatar,
  MenuItem,
  Typography,
  Toolbar,
  AppBar,
  Button,
  ClickAwayListener,
  Grow,
  MenuList,
  Paper,
  Popper,
  Switch,
} from "@material-ui/core";

// auth
import { useAuth0 } from "@auth0/auth0-react";

const useStyles = makeStyles((theme) => ({
  grow: {
    flexGrow: 1,
  },
  title: {
    fontSize: "2rem",
  },
  small: {
    width: theme.spacing(3),
    height: theme.spacing(3),
    marginRight: "5px",
  },
  userButton: {
    textTransform: "none",
  },
}));

export default function MenuBar({ darkThemeIcon, darkTheme, setDarkTheme }) {
  const classes = useStyles();
  const [open, setOpen] = useState(false);
  const anchorRef = useRef(null);

  const [loggedInUser, setLoggedInUser] = useState({ name: "", picUrl: "" });

  const {
    isAuthenticated,
    logout,
    loginWithRedirect,
    user,
    getIdTokenClaims,
  } = useAuth0();

  const handleToggle = () => {
    getIdTokenClaims().then((data) => {
      console.log(data);
    });
    console.log(user);
    setOpen((prevOpen) => !prevOpen);
  };

  const handleClose = (event) => {
    if (anchorRef.current && anchorRef.current.contains(event.target)) {
      return;
    }

    setOpen(false);
  };

  const handleLogout = () => {
    setOpen(false);
    logout({ returnTo: window.location.origin });
  };

  const handleLogin = () => {
    loginWithRedirect();
  };

  const prevOpen = useRef(open);
  useEffect(() => {
    if (prevOpen.current === true && open === false) {
      anchorRef.current.focus();
    }

    prevOpen.current = open;
  }, [open]);

  useEffect(() => {
    if (user) {
      setLoggedInUser({ name: user.given_name, picUrl: user.picture });
    }
  }, [user]);

  const renderMenu = (
    <Popper
      style={{ zIndex: "2" }}
      open={open}
      anchorEl={anchorRef.current}
      role={undefined}
      transition
      disablePortal
    >
      {({ TransitionProps, placement }) => (
        <Grow
          {...TransitionProps}
          style={{
            transformOrigin:
              placement === "bottom" ? "center top" : "center bottom",
          }}
        >
          <Paper>
            <ClickAwayListener onClickAway={handleClose}>
              <MenuList autoFocusItem={open} id="menu-list-grow">
                <MenuItem onClick={handleClose}>Profile</MenuItem>
                <MenuItem onClick={handleClose}>My account</MenuItem>
                <MenuItem onClick={handleLogout}>Logout</MenuItem>
              </MenuList>
            </ClickAwayListener>
          </Paper>
        </Grow>
      )}
    </Popper>
  );

  return (
    <div className={classes.grow}>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h5" noWrap>
            Couch Potatoes
          </Typography>
          <div className={classes.grow} />
          <div>
            {darkThemeIcon}
            <Switch
              checked={darkTheme}
              onChange={() => {
                setDarkTheme(!darkTheme);
              }}
            />
            {isAuthenticated ? (
              <>
                <Button
                  className={classes.userButton}
                  aria-label="account of current user"
                  ref={anchorRef}
                  aria-controls={open ? "menu-list-grow" : undefined}
                  aria-haspopup="true"
                  onClick={handleToggle}
                  color="inherit"
                  // startIcon={<AccountCircle />}
                >
                  <Avatar className={classes.small} src={`${loggedInUser.picUrl}`} />
                  {loggedInUser.name}
                </Button>
              </>
            ) : (
              <Button onClick={handleLogin} color="default" variant="contained">
                Log in
              </Button>
            )}
          </div>
        </Toolbar>
      </AppBar>
      {renderMenu}
    </div>
  );
}
