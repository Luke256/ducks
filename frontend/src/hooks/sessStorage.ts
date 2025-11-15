'use client';

import { useEffect, useState } from "react";


const useSessionStorage = <T,>(key: string, initialValue: T) => {
    const [storedValue, setStoredValue] = useState<T>(initialValue);
    useEffect(() => {
        try {
            const item = window.sessionStorage.getItem(key);
            setStoredValue(item ? JSON.parse(item) : initialValue);
        } catch (error) {
            setStoredValue(initialValue);
        }
    }, []);

    const setValue = (value: T | ((val: T) => T)) => {
        try {
            const valueToStore = value instanceof Function ? value(storedValue) : value;
            setStoredValue(valueToStore);
            window.sessionStorage.setItem(key, JSON.stringify(valueToStore));
        } catch (error) { }
    };

    return [storedValue, setValue] as const;
}

export { useSessionStorage };