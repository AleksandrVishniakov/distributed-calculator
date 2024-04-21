import React, {useState} from "react";
import ErrorMessage from "../ErrorMessage/ErrorMessage";
import useSavedState from "../../../hooks/useSavedState";

import './AppRegisterForm.css'

interface AppRegisterFormProps {
    name: string
    changeScreensTitle?: string

    onChangeScreens?: () => void
    onSubmit?: (login: string, password: string) => void
    onBack?: () => void

    showError?: boolean
    errorTitle?: string
    errorMessage?: string
}

const AppRegisterForm: React.FC<AppRegisterFormProps> = (props) => {
    const [emailInputValue, setEmailInputValue] = useSavedState("", "login-input-value")
    const [passwordInputValue, setPasswordInputValue] = useState("")

    const [isDataSubmitting, setDataSubmitting] = useState<boolean>(false)

    const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        setDataSubmitting(true)

        try {
            if (props.onSubmit) {
                props.onSubmit(emailInputValue, passwordInputValue)
            }
        } catch {
        } finally {
            setDataSubmitting(false)
        }
    }

    const handleBack = () => {
        if (props.onBack) {
            props.onBack()
        }
    }

    const handleChangeScreens = () => {
        if (props.onChangeScreens) {
            props.onChangeScreens()
        }
    }

    const onInputChange = (callback: React.Dispatch<React.SetStateAction<string>>) => {
        return (evt: React.ChangeEvent<HTMLInputElement>) => {
            callback(evt.target.value)
        }
    }

    return (
        <form
            className="AppRegisterForm_form"
            onSubmit={onSubmit}
        >
            <h2 className="AppRegisterForm_title">{props.name}</h2>

            <input
                type="text"
                placeholder="login"
                name="login"
                value={emailInputValue}
                onChange={onInputChange(setEmailInputValue)}
                disabled={isDataSubmitting}
            />
            <input
                type="password"
                name="password"
                placeholder="password"
                onChange={onInputChange(setPasswordInputValue)}
                value={passwordInputValue}
                disabled={isDataSubmitting}
            />

            {props.changeScreensTitle ?
                <p
                    onClick={handleChangeScreens}
                    className="custom-link">
                    {props.changeScreensTitle}
                </p> : null}

            <button
                className="AppRegisterForm_submit-btn"
                type="submit"
                disabled={isDataSubmitting}
            >
                {props.name}
            </button>

            {props.showError ? <ErrorMessage
                title={props.errorTitle ? props.errorTitle : ""}
                message={props.errorMessage ? props.errorMessage : ""}
            /> : null}

            <p
                className="custom-link back-btn"
                onClick={handleBack}
            >
                На главную
            </p>
        </form>
    )
}

export default AppRegisterForm