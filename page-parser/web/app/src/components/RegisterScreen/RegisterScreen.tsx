import React from "react";
import AppTitle from "../common/AppTitle/AppTitle";
import SidePanel from "../common/SidePanel/SidePanel";
import AppRegisterForm from "../common/AppRegisterForm/AppRegisterForm";

import "./RegisterScreen.css"
import {AuthAPI} from "../../api/AuthAPI";

interface RegisterScreenProps {
    authAPI: AuthAPI

    onChangeRegisterScreens?: ()=>void
    onHome: ()=>void
    onError: (msg: string)=>void
    onSuccess: ()=>void
}

const RegisterScreen: React.FC<RegisterScreenProps> = (props) => {
    const handleScreenChange = () => {
        if (props.onChangeRegisterScreens) {
            props.onChangeRegisterScreens()
        }
    }

    const handleSubmit = (login: string, password: string) => {
        props.authAPI.register(login, password)
            .then(props.onSuccess)
            .catch((e) => {
                props.onError(e.toString())
            })
    }

    return (
        <section className="RegisterScreen">
            <SidePanel>
                <>
                    <AppTitle/>

                    <AppRegisterForm
                        name="Регистрация"
                        changeScreensTitle="Есть аккаунт? Вход"
                        onChangeScreens={handleScreenChange}
                        onSubmit={handleSubmit}
                        onBack={props.onHome}
                    />

                    <div></div>
                </>
            </SidePanel>

            <div className="background">

            </div>
        </section>
    )
}

export default RegisterScreen