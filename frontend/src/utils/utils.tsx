import moment from "moment";
import { DATE_FORMAT } from "./const";
export const stringifyJSON = function(obj) {
    if (typeof obj === "number" || obj === null || typeof obj === "boolean") {
        return `${obj}`;
    } else if (typeof obj === "string") {
        if (
            (obj.charAt(0) === "{" && obj.charAt(obj.length - 1) === "}") ||
            (obj.charAt(0) === "[" && obj.charAt(obj.length - 1) === "]")
        ) {
            return stringifyJSON(JSON.parse(obj));
        }
        return `"${obj}"`;
    }
    let bin = [];
    if (Array.isArray(obj)) {
        if (obj.length === 0) {
            return "[]";
        } else {
            for (let i = 0; i < obj.length; i++) {
                if (typeof obj[i] === "string") {
                    let str = stringifyJSON(obj[i]);
                    bin.push(str);
                } else if (Array.isArray(obj[i])) {
                    let arr = stringifyJSON(obj[i]);
                    bin.push(arr);
                } else if (typeof obj[i] === "number") {
                    bin.push(obj[i]);
                } else {
                    let ifObj = stringifyJSON(obj[i]);
                    bin.push(ifObj);
                }
            }
            return `[${bin}]`;
        }
    } else {
        let createdArr = [];
        if (Object.keys(obj).length === 0) {
            return "{}";
        } else {
            for (let key in obj) {
                if (
                    typeof obj[key] === "string" ||
                    typeof obj[key] === "boolean" ||
                    obj[key] === null
                ) {
                    let strKey = stringifyJSON(key);
                    let strVal = stringifyJSON(obj[key]);
                    let strArr = strKey + ":" + strVal;
                    createdArr.push(strArr);
                } else if (Array.isArray(obj[key])) {
                    let arrKey = stringifyJSON(key);
                    let arrVal = stringifyJSON(obj[key]);
                    let arrArr = arrKey + ":" + arrVal;
                    createdArr.push(arrArr);
                } else if (
                    typeof obj[key] === "function" ||
                    obj[key] === undefined
                ) {
                    delete obj[key];
                    stringifyJSON(obj);
                } else {
                    let objKey = stringifyJSON(key);
                    let objVal = stringifyJSON(obj[key]);
                    let objObj = objKey + ":" + objVal;
                    createdArr.push(objObj);
                }
            }
        }
        return `{${createdArr}}`;
    }
};

export const formatNumberWithComma = x => {
    if (!x) {
        return 0;
    }
    let parts = x.toString().split(".");
    parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ",");
    return parts.join(".");
};

export const makeSortList = arr => {
    return [].concat(arr).sort((a, b) => a.name.localeCompare(b.name));
};

export const encodeURIMoment = _moment => {
    if (!_moment) {
        return;
    }
    return encodeURIComponent(
        moment
            .utc(_moment)
            .add(9, "hour")
            .format("YYYY-MM-DDTHH:mm:00+09:00")
    );
};

export const formatMoment = _moment => {
    if (!_moment) {
        return;
    }
    return moment(_moment).format(DATE_FORMAT);
};
