import React, {PropsWithChildren} from "react";
import './Header.module.css'

const Header: React.FC<PropsWithChildren> = (props) => {
    return (
        <header>
            {props.children}
        </header>
    )
}

export default Header