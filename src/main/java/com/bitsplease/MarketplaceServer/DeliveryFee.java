package com.bitsplease.MarketplaceServer;

import javax.jws.WebMethod;
import javax.jws.WebService;

@WebService()
public class DeliveryFee {
    @WebMethod
    public double CalculateDelivery(double distance) {
        double result = distance * 0.2;
        System.out.println("Delivery fee: " + result);
        return result;
    }
}
