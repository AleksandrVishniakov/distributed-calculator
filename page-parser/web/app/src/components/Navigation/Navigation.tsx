import React, {useState} from "react";
import './Navigation.module.css'

type NavigationElement = [string, ()=>void]

interface NavigationProps {
    elements: NavigationElement[]
}

const Navigation: React.FC<NavigationProps> = (props) => {
    const [navIndex, setNavIndex] = useState(0)

    const handleNavClick = (onClick: ()=>void, index: number)=> {
        if (index === navIndex) return

        setNavIndex(index)
        onClick()
    }

    return (
        <nav>
            <ol>
                {
                    props.elements.map((element, index) => {
                        return (
                            <li
                                //id={`nav-el-${index}`}
                                key={index}
                                onClick={()=>{
                                    handleNavClick(element[1], index)
                                }}

                                className={
                                    index === navIndex ? "chosen" : ""
                                }
                            >
                                {element[0]}
                            </li>
                        )
                    })
                }
            </ol>
        </nav>
    )
}

export default Navigation