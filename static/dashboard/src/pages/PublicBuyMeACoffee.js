import { Component } from "react";
import CreatorPage from "./CreatorPage";
import CreatorRecentSupporters from "../components/CreatorRecentSupporters";
import PaymentUI from "../components/payments/PaymentUI";


class PublicBuyMeACoffee extends Component{
    render(){
        return (
            <div class="padding-bottom-40 padding-top-10" >
                    <div class="row justify-content-center" >
                        <div class="col-lg-7 col-xs-12 col-md-7 order-md-1">
                            <PaymentUI purpose='support' user={this.props.user} amount={this.props.amount}  items={this.props.items} itemName={this.props.itemName} creator={this.props.creator.username}  />
                        </div>
                        <div class="col-lg-5 col-xs-12 col-md-5 order-md-2">
                            <CreatorRecentSupporters grandMax={5} maxShowMobile={2} supporters={this.props.supporters} creator={this.props.creator} user={this.props.user} />
                        </div>
                        
                    </div>
                    
                </div>
        )
    }
}

export default PublicBuyMeACoffee;