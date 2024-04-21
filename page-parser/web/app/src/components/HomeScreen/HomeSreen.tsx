import React from "react";
import AppTitle from "../common/AppTitle/AppTitle";

import './HomeScreen.css'

interface Props {
    onLogin : ()=>void
    onRegister : ()=>void
}

const HomeScreen: React.FC<Props> = ({onLogin, onRegister}) => {
    return (
        <section className="HomeScreen">
            <AppTitle/>
            <p>Распределённый калькулятор - приложение для подсчёта арифметичсекого выражения на разных машинах. Вы можете задать время выполнения каждой операции, а также увидеть, задачи каждого агента</p>

            <div className="HomeScreen__actions">
                <button
                    onClick={onLogin}
                >Войдите</button>
                или
                <button
                    onClick={onRegister}
                >Зарегистрируйтесь</button>
            </div>
        </section>
    )
}

export default HomeScreen