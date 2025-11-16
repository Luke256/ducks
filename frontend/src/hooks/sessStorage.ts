'use client';

import { useCallback, useSyncExternalStore } from "react";


const useSessionStorage = <T,>(key: string, initialValue: T) => {
    const subscribe = useCallback((callback: () => void) => {
        if (typeof window === "undefined") {
            return () => { };
        }

        const handler = (event: Event) => {
            if (event instanceof StorageEvent && event.key !== key) {
                return;
            }
            callback();
        };

        window.addEventListener("storage", handler);
        window.addEventListener("session-storage", handler);

        return () => {
            window.removeEventListener("storage", handler);
            window.removeEventListener("session-storage", handler);
        };
    }, [key]);

    const getSnapshot = useCallback(() => {
        if (typeof window === "undefined") {
            return initialValue;
        }
        try {
            const item = window.sessionStorage.getItem(key);
            return item ? JSON.parse(item) : initialValue;
        } catch {
            return initialValue;
        }
    }, [key, initialValue]);

    const getServerSnapshot = useCallback(() => initialValue, [initialValue]);

    const storedValue = useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);

    const setValue = (value: T | ((val: T) => T)) => {
        if (typeof window === "undefined") {
            return;
        }

        try {
            const valueToStore = value instanceof Function ? value(getSnapshot()) : value;
            window.sessionStorage.setItem(key, JSON.stringify(valueToStore));
            window.dispatchEvent(new Event("session-storage"));
        } catch { }
    };

    return [storedValue, setValue] as const;
}

export { useSessionStorage };