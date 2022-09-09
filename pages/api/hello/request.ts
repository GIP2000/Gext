/* Do not change, this code is generated from Golang structs */


export interface Something {
    Name: string;
    Value: number;
}
export const getData = async ():Promise<Something> =>(await fetch("/api/hello")).json() as Promise<Something>