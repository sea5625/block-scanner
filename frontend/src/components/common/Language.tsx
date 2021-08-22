import { createContext, useState, useEffect } from "react";
import i18n from "utils/i18n";

const Language = () => {
    const globalLanguage =
        sessionStorage.getItem("redux") &&
        JSON.parse(sessionStorage.getItem("redux")).storage.language
            ? JSON.parse(sessionStorage.getItem("redux")).storage.language
            : false;

    const [language, setLanguage] = useState(
        globalLanguage ? globalLanguage : "ko"
    );
    useEffect(() => {
        setLanguage(globalLanguage);
    }, [globalLanguage]);

    const setGlobalLanguage = (lang: string) => {
        setLanguage(lang);
    };

    const I18n = i18n[language];
    const LanguageContext = createContext({ I18n, setGlobalLanguage });
    return LanguageContext;
};

export default Language;
