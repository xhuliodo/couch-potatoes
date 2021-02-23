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
import "./MenuBar.css";

// auth
import { useAuth0 } from "@auth0/auth0-react";
import { Brightness3, Brightness7 } from "@material-ui/icons";

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
    marginRight: "7px",
  },
  userButton: {
    textTransform: "none",
  },
}));

export default function MenuBar({ darkThemeIcon, darkTheme, setDarkTheme }) {
  const classes = useStyles();
  const [open, setOpen] = useState(false);
  const anchorRef = useRef(null);

  const {
    isAuthenticated,
    logout,
    loginWithRedirect,
    user,
    getIdTokenClaims,
  } = useAuth0();

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

  const handleLogout = () => {
    localStorage.removeItem("finishedSetup");
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
                {/* <MenuItem onClick={handleClose}>Profile</MenuItem> */}
                {/* <MenuItem onClick={handleClose}>My account</MenuItem> */}
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
        <Toolbar style={{ minHeight: "48px" }}>
          <Typography variant="h5" noWrap>
            Couch Potatoes
          </Typography>
          <div className={classes.grow} />
          <div>
            <Switch
              checked={darkTheme}
              icon={<Brightness7 />}
              checkedIcon={<Brightness3 />}
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
                  // startIcon={<AccountCircle />}
                >
                  <Avatar
                    className={classes.small}
                    src={user ? user.picture : ""}
                  />
                  {user.given_name}
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
