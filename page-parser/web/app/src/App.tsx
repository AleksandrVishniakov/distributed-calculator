import React, {forwardRef, useState} from 'react';
import './App.css';
import useSavedState from "./hooks/useSavedState";
import RegisterScreen from "./components/RegisterScreen/RegisterScreen";
import LoginScreen from "./components/LoginScreen/LoginScreen";
import ActionsScreen from "./components/ActionsScreen/ActionsScreen";
import HomeScreen from "./components/HomeScreen/HomeSreen";
import {AuthAPI} from "./api/AuthAPI";
import {ExpressionsAPI} from "./api/ExpressionsAPI";
import {OperatorsAPI} from "./api/OperatorsAPI";
import {WorkerAPI} from "./api/WorkersAPI";
import MuiAlert, { AlertProps } from '@mui/material/Alert';
import Snackbar from '@mui/material/Snackbar';

enum Screens {
    Home,
    Register,
    Login,
    Actions
}

function AppNavigation(
    authAPI : AuthAPI,
    expressionsAPI : ExpressionsAPI,
    operatorsAPI : OperatorsAPI,
    workersAPI : WorkerAPI,
    screen: Screens,
    setScreen: React.Dispatch<React.SetStateAction<number>>,
    onError: (msg: string)=>void
): React.ReactNode {
    switch (screen) {
        case Screens.Register:
            return (
                <RegisterScreen
                    authAPI={authAPI}

                    onChangeRegisterScreens={()=>{
                        setScreen(Screens.Login)
                    }}

                    onHome={()=>{
                        setScreen(Screens.Home)
                    }}

                    onError={onError}
                    onSuccess={()=>{
                        setScreen(Screens.Login)
                    }}
                />
            )

        case Screens.Login:
            return (
                <LoginScreen
                    authAPI={authAPI}


                    onChangeRegisterScreens={()=>{
                        setScreen(Screens.Register)
                    }}

                    onHome={()=>{
                        setScreen(Screens.Home)
                    }}

                    onError={onError}
                    onSuccess={()=>{
                        setScreen(Screens.Actions)
                    }}
                />
            )

        case Screens.Actions:
            return (
                <ActionsScreen
                    expressionsAPI={expressionsAPI}
                    operatorsAPI={operatorsAPI}
                    workersAPI={workersAPI}

                    user = {{
                        login: "выйти"
                    }}

                    onLogout={()=>{
                        setScreen(Screens.Home)
                        localStorage.removeItem("token")
                    }}
                    onError={onError}
                />
            )

        case Screens.Home:
            return (
                <HomeScreen
                    onLogin={()=>{
                        setScreen(Screens.Login)
                    }}
                    onRegister={()=>{
                        setScreen(Screens.Register)
                    }}
                />
            )
    }

}

const App: React.FC<{
    authAPI: AuthAPI
    expressionsAPI: ExpressionsAPI
    operatorsAPI: OperatorsAPI
    workersAPI : WorkerAPI
}> = (props) => {
    const [screen, setScreen] = useSavedState<number>(Screens.Actions, "screen")
    const [theme, setTheme] = useSavedState<"light" | "dark">("dark", "theme")
    const [errorSnackbarOpen, setErrorSnackbarOpen] = useState(false)
    const [errorSnackbarText, setErrorSnackbarText] = useState("")

    const toggleTheme = () => {
        setTheme(theme === "light" ? "dark" : "light");
    }

    const handleSnackbarClose = (event?: React.SyntheticEvent | Event, reason?: string) => {
        if (reason === 'clickaway') {
            return;
        }

        setErrorSnackbarOpen(false);
    };

    const handleError = (msg: string) => {
        if (msg.split(": ")[1] === "401") {
            setScreen(Screens.Home)
        }
        console.error(msg)

        setErrorSnackbarText(msg)
        setErrorSnackbarOpen(true)
    }

    return (
        <div className={`App ${theme}`}>
            {
                AppNavigation(
                    props.authAPI,
                    props.expressionsAPI,
                    props.operatorsAPI,
                    props.workersAPI,
                    screen,
                    setScreen,
                    handleError
                )
            }

            <Snackbar open={errorSnackbarOpen} autoHideDuration={5000} onClose={handleSnackbarClose}
                      anchorOrigin={{vertical: 'bottom', horizontal: 'center'}}>
                <Alert onClose={handleSnackbarClose} severity="error" sx={{width: '100%'}}>
                    {errorSnackbarText}
                </Alert>
            </Snackbar>
        </div>
    );
}

const Alert = forwardRef<HTMLDivElement, AlertProps>(function Alert(
    props,
    ref,
) {
    return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
});

export default App;
