/* Do not change, this code is generated from Golang structs */


export interface BR {
    Status: number;
}
export const getData = async ():Promise<BR> =>(await fetch("/api/badRequest")).json() as Promise<BR>
