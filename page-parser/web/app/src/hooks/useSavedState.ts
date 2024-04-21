import React, { useEffect, useState } from "react";

const useSavedState = <T>(defaultState: T, name: string): [T, React.Dispatch<React.SetStateAction<T>>] => {
    const getInitialValue = (): T => {
        const savedState = window.localStorage.getItem(name);
        if (savedState) {
            try {
                return JSON.parse(savedState) as T;
            } catch (error) {
                console.error("useSavedState: error parsing saved state:", error);
            }
        }
        return defaultState;
    };

    const [state, setState] = useState<T>(getInitialValue());

    useEffect(() => {
        window.localStorage.setItem(name, JSON.stringify(state));
    }, [name, state]);

    return [state, setState];
};

export default useSavedState;