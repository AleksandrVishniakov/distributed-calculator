import React, {forwardRef, useState} from 'react';
import './App.css';
import Header from "./components/Header/Header";
import Navigation from "./components/Navigation/Navigation";
import CalculationScreen from "./components/CalculationsScreen/CalculationScreen";
import SettingsScreen from "./components/SettingsScreen/SettingsScreen";
import WorkersScreen from "./components/WorkersScreen/WorkersScreen";
import MuiAlert, { AlertProps } from '@mui/material/Alert';
import Snackbar from '@mui/material/Snackbar';


enum Screens {
    Calculation,
    Workers,
    Settings,
}

interface AppProps {
    host: string
}

const App: React.FC<AppProps> = (props) => {
    const [currentScreen, setScreen] = useState(Screens.Calculation)
    const [errorSnackbarOpen, setErrorSnackbarOpen] = useState(false)
    const [errorSnackbarText, setErrorSnackbarText] = useState("")

    const handleSnackbarClose = (event?: React.SyntheticEvent | Event, reason?: string) => {
        if (reason === 'clickaway') {
            return;
        }

        setErrorSnackbarOpen(false);
    };

    const handleAppError = (msg: string) => {
        console.log(msg)
        setErrorSnackbarText(msg)
        setErrorSnackbarOpen(true)
    }

    const renderScreen = () => {
        switch (currentScreen) {
            case Screens.Calculation:
                return (
                    <CalculationScreen
                        onClick={(expr) => {

                        }}

                        host={props.host}
                        onError={handleAppError}
                    />
                )

            case Screens.Settings:
                return <SettingsScreen host={props.host} onError={handleAppError}/>
            case Screens.Workers:
                return <WorkersScreen host={props.host} onError={handleAppError}/>
        }
    }

    return (
        <div className="App">
            <div style={{position: "sticky", width: "100%"}}>
                <Header>
                    <h1>Распределённый калькулятор</h1>
                </Header>

                <Navigation elements={[
                    ["посчитать выражение", () => {
                        setScreen(Screens.Calculation)
                    }],
                    ["доступные машины", () => {
                        setScreen(Screens.Workers)
                    }],
                    ["настройки", () => {
                        setScreen(Screens.Settings)
                    }],
                ]}/>
            </div>

            <main>
                {
                    renderScreen()
                }
            </main>

            <Snackbar open={errorSnackbarOpen} autoHideDuration={5000} onClose={handleSnackbarClose}
                      anchorOrigin={{vertical: 'bottom', horizontal: 'center'}}>
                <Alert onClose={handleSnackbarClose} severity="error" sx={{width: '100%'}}>
                    {errorSnackbarText}
                </Alert>
            </Snackbar>
        </div>
);
}

export default App;

const Alert = forwardRef<HTMLDivElement, AlertProps>(function Alert(
    props,
    ref,
) {
    return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
});