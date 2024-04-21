import React from "react";
import AppIcon from "./AppIcon/AppIcon";

import "./AppTitle.css";

const AppTitle: React.FC<{ onClick?: () => void }> = ({onClick}) => {
    return (
        <div className="AppTitle" onClick={onClick}>
            <AppIcon size="large"/>

            <h1>Распределённый калькулятор</h1>
        </div>
    )
}

export default AppTitle