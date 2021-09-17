import os from 'os'
import BigNumber from 'bignumber.js'
let Bech32 = require('bech32')
let Sha256 = require("sha256")
export function getIpAddress() {
    const interfaces = os.networkInterfaces();
    for (const devName in interfaces) {
        const iface = interfaces[devName];
        for (let i = 0; i < iface.length; i++) {
            let alias = iface[i];
            if (alias.family === 'IPv4' && alias.address !== '127.0.0.1' && !alias.internal) {
                return alias.address;
            }
        }
    }
}

export function getTimestamp(): number {
    return Math.floor(new Date().getTime() / 1000);
}

export function formatDateStringToNumber(dateString) {
    return Math.floor(new Date(dateString).getTime() / 1000)
}

export function addressTransform(str:string, prefix?:string) {
    try {
        let Bech32str = Bech32.decode(str,'utf-8');
        prefix = prefix || '';
        let result = Bech32.encode(prefix, Bech32str.words)
        return result;
    } catch (e) {
        console.warn('address transform failed', e)
    }
}

export function hexToBech32(hexStr:string, prefix:string = "") {
    try {
        let words = Bech32.toWords(Buffer.from(hexStr,'hex'));
        return Bech32.encode(prefix, words);
    }catch (e) {
        console.warn('address transform fialed',e)
    }
}

export function pageNation(dataArray: any[], pageSize: number = 0) {
    let index: number = 0;
    let newArray: any = [];
    if (dataArray.length > pageSize) {
        while (index < dataArray.length) {
            newArray.push(dataArray.slice(index, index += pageSize));
        }
    } else {
        newArray = dataArray
    }
    return newArray
}

export function getAddress(publicKey) {
    let words = Bech32.decode(publicKey).words;
    words =  Bech32.fromWords(words);
    if (words.length > 33){
        //去掉amino编码前缀
        words = words.slice(5)
    }
    let addr = Sha256(Buffer.from(words));
    if (addr && addr.length > 40) {
        addr = addr.substr(0,40);
    }
    return addr;
}

export function splitString(str,symbol) {
    let array = str.split(symbol)
    return array[array.length - 1]
}

export function uniqueArr(arr, brr) {
    let temp:string[] = [];
    let temparray:string[] = [];
    for (var i = 0; i < brr.length; i++) {  
        temp[brr[i]] = typeof brr[i];;
    }
    for (var i = 0; i < arr.length; i++) {  
        var type = typeof arr[i];
        if (!temp[arr[i]]) {  
            temparray.push(arr[i]);
        } else if (temp[arr[i]].indexOf(type) < 0) { 
            temparray.push(arr[i]); 
        }  
    }
    return temparray
}

export function BigNumberPlus(num1, num2) {
    const x = new BigNumber(num1)
    return x.plus(num2).toNumber()
}

export function BigNumberDivision(num1, num2) {
    const x = new BigNumber(num1)
    const y = new BigNumber(num2)
    return x.dividedBy(y).toFixed()  
}

export function BigNumberMinus(num1, num2) {
    const x = new BigNumber(num1)
    return x.minus(num2).toNumber()  
}

export function BigNumberMultiplied(num1, num2) {
    const x = new BigNumber(num1)
    return x.multipliedBy(num2).toNumber()  
}
