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

// mendja eshte gje e madhe, tek file per 2 px XD
import "./MenuBar.scss";

// auth
import { useAuth0 } from "@auth0/auth0-react";
import { authConfig } from "../utils/auth0";

// dark mode icons
import { Brightness3, Brightness7 } from "@material-ui/icons";
// conditional landing page styling import
import { useHistory, useLocation } from "react-router-dom";

const useStyles = makeStyles((theme) => ({
  grow: {
    flexGrow: 1,
    maxHeight: "50px",
  },
  title: {
    fontSize: "1.3rem",
    cursor: "pointer",
  },
  small: {
    width: theme.spacing(3),
    height: theme.spacing(3),
    marginRight: "7px",
  },
  userButton: {
    textTransform: "none",
    color: "#fff",
  },
}));

export default function MenuBar({ darkTheme, setDarkTheme }) {
  const classes = useStyles();

  // menu bar in landing page show/hide
  const location = useLocation();
  const history = useHistory();
  const [isUserInLandingPage, setIsUserInLandingPage] = useState(true);

  useEffect(() => {
    if (location.pathname === "/") {
      setIsUserInLandingPage(true);
    } else setIsUserInLandingPage(false);
  }, [location]);

  // auth logic
  const {
    isAuthenticated,
    logout,
    loginWithRedirect,
    user,
    getIdTokenClaims,
  } = useAuth0();

  const handleLogout = () => {
    localStorage.removeItem("finishedSetup");
    setOpen(false);
    logout({ client_id: authConfig.client_id, returnTo: authConfig.logoutUri });
  };

  const handleLogin = () => {
    loginWithRedirect();
  };
  // toolbar logic
  const [open, setOpen] = useState(false);
  const anchorRef = useRef(null);

  const handleToggle = async () => {
    setOpen((prevOpen) => !prevOpen);
    const token = await getIdTokenClaims();
    console.log(token.__raw);
    console.log(user);
  };

  const handleClose = (event) => {
    if (anchorRef.current && anchorRef.current.contains(event.target)) {
      return;
    }

    setOpen(false);
  };

  const prevOpen = useRef(open);
  useEffect(() => {
    if (prevOpen.current === true && open === false) {
      anchorRef.current.focus();
    }

    prevOpen.current = open;
  }, [open]);

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
          <Paper elevation={2}>
            <ClickAwayListener onClickAway={handleClose}>
              <MenuList autoFocusItem={open} id="menu-list-grow">
                <MenuItem onClick={handleLogout}>Logout</MenuItem>
              </MenuList>
            </ClickAwayListener>
          </Paper>
        </Grow>
      )}
    </Popper>
  );

  return isUserInLandingPage ? null : (
    <div className={classes.grow}>
      <AppBar position="static">
        <Toolbar style={{ minHeight: "48px" }}>
          <Typography
            noWrap
            className={classes.title}
            onClick={() => {
              history.push("/");
            }}
          >
            Couch Potatoes
          </Typography>
          <div className={classes.grow} />
          <div>
            <Switch
              checked={darkTheme}
              icon={<Brightness7 />}
              color="default"
              checkedIcon={<Brightness3 style={{ color: "#f2ebff " }} />}
              onChange={() => {
                localStorage.setItem("darkMode", !darkTheme);
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
                >
                  <Avatar
                    className={classes.small}
                    src={user ? user?.picture : ""}
                  />
                  {user?.given_name}
                </Button>
              </>
            ) : (
              <Button
                className={classes.userButton}
                variant="outlined"
                onClick={handleLogin}
              >
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
