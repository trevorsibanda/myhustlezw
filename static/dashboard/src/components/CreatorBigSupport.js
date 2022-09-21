import React, { Component } from "react"
import { Link, Switch, Route } from 'react-router-dom'

function CreatorBigSupport(props) {
    return (
<div class="single-price-plan-01 d-none d-md-block">
    <div class="price-header">
        <h4 class="title">Buy { props.creator.fullname } a { props.creator.page.donation_item }</h4>
    </div>
    <div class="price-wrap">
        <span class="price">{ props.creator.page.donation_item_unit_price }</span>
    </div>
    
    <div class="price-body">
        <ul>
            <li>
                <label>Your name</label>
                <input type="text" class="form-control" id="name" placeholder="Your name (optional)" />
            </li>
            <li>
                <label>* Your Ecocash phone number</label>
                <input type="tel" placeholder="07XXXXXXXX (required)" required id="phone" class="form-control" />
            </li>
            <li>
                <small>You can leave a message after you've completed the payment.</small>
            </li>
            
        </ul>
        
        
    </div>
    <div class="price-footer">
        <div class="btn-wrapper">
            <button class="boxed-btn" onclick="payEcocash();" >Pay with Ecocash <strong>{ props.creator.page.donation_item_unit_price }</strong></button>
        </div>
        <small>All Ecocash payments are handled through <a href="https://paynow.co.zw/" rel="noreferrer" target="_blank">PayNow</a></small>
        <p><strong>OR</strong></p>
        <small> <a href="#" class="btn btn-block btn-sm btn-info">Pay with 2Checkout (International Mastercard/Skrill)</a></small>
    </div>
</div>
    )
}

export default CreatorBigSupport;