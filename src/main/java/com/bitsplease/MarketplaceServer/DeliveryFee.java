package com.bitsplease.MarketplaceServer;

import javax.jws.WebMethod;
import javax.jws.WebService;

@WebService()
public class DeliveryFee {
    @WebMethod
    public double CalculateDelivery(double weight) {
        double result = DeliveryFeeProcess(weight);
        System.out.println("Delivery fee: " + result);
        return result;
    }

    public static double DeliveryFeeProcess(double weight) {
        return weight * 0.2;
    }
}
