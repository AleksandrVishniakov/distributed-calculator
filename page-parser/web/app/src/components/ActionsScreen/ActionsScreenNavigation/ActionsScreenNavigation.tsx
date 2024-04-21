import React, {useEffect} from "react";
import useSavedState from "../../../hooks/useSavedState";

import './ActionsScreenNavigation.css'

interface NavItemProps {
    title: string
    Icon: React.ReactNode

    onClick?: () => void
}

interface Props {
    defaultItem: string
    Items: Array<NavItemProps>;
}

const ActionsScreenNavigation: React.FC<Props> = ({ defaultItem, Items}) => {
    const [currentItem, setCurrentItem] = useSavedState(defaultItem, "actions-screen-nav")

    useEffect(() => {
        Items.forEach((value)=> {
            if (value.title === currentItem) {
                if (value.onClick) {
                    value.onClick()
                    return
                }
            }
        })
    }, [currentItem, Items]);

    const handleItemClick = (title: string) => {
        setCurrentItem(title)
    }

    return (
        <nav className="ActionsScreenNavigation">
            {Items.map((item) => {
                return (
                    <div
                        className={`ActionsScreenNavigation__item ${item.title === currentItem ? "selected" : ""}`}
                        onClick={()=>{handleItemClick(item.title)}}
                        key={item.title}
                        title={item.title}
                    >
                        {item.Icon}
                        <p className="ActionsScreenNavigation__title">
                            {item.title}
                        </p>
                    </div>
                )
            })}
        </nav>
    )
}

export default ActionsScreenNavigation;