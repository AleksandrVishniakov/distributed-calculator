import React, {PropsWithChildren} from "react";

import "./SidePanel.css"

const SidePanel: React.FC<PropsWithChildren> = ({children}) => {
    return (
        <section className="SidePanel">
            {children}
        </section>
    )
}

export default SidePanel