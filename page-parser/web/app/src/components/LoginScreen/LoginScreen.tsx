import SidePanel from "../common/SidePanel/SidePanel";
import AppTitle from "../common/AppTitle/AppTitle";
import AppRegisterForm from "../common/AppRegisterForm/AppRegisterForm";
import React from "react";

import './LoginScreen.css'
import {AuthAPI} from "../../api/AuthAPI";

interface LoginScreenProps {
    authAPI: AuthAPI

    onChangeRegisterScreens?: ()=>void
    onHome: ()=>void
    onError: (msg: string)=>void
    onSuccess: ()=>void
}

const LoginScreen: React.FC<LoginScreenProps> = (props) => {
    const handleScreenChange = () => {
        if (props.onChangeRegisterScreens) {
            props.onChangeRegisterScreens()
        }
    }

    const handleSubmit = (login: string, password: string) => {
        props.authAPI.login(login, password)
            .then(props.onSuccess)
            .catch((e) => {
                props.onError(e.toString())
            })
    }

    return (
        <section className="LoginScreen">
            <SidePanel>
                <>
                    <AppTitle/>

                    <AppRegisterForm
                        name="Вход"
                        changeScreensTitle="Нет аккаунта? Регистрация"
                        onChangeScreens={handleScreenChange}
                        onBack={props.onHome}
                        onSubmit={handleSubmit}
                    />

                    <div></div>
                </>
            </SidePanel>

            <div className="background">

            </div>
        </section>
    )
}

export default LoginScreen