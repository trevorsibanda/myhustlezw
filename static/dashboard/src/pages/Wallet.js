import React, { Component, useState } from "react";
import { Link, Redirect } from "react-router-dom"
import v1 from "../api/v1";

import WalletSummary from "../components/WalletSummary"
import RecentWalletOperations from "../components/RecentWalletActivity"
import WalletRequestWithdrawal from "../components/WalletRequestWithdrawal"
import WalletUpdatePayoutDetails from "../components/WalletPayoutDetails"
import ActionCenter from "../components/ActionCenter";
import WalletAPISettings from "../components/WalletAPISettings";



class Wallet extends Component {
    constructor(props) {

        v1.page.set({ title: 'My Wallet / Dashboard' })
        v1.page.track()
        
        super(props)

        this.state = {
            currency: "ZWL",
            summary: {
                "USD":  {
                    "currency": "USD",
                    "escrow": 0,
                    "available": 0,
                    "withdrawn": 0,
                    "disputed": 0,
                    "refunded": 0,
                    "inaccurate": false
                }, 
                "ZWL": { 
                    "currency": "ZWL",
                    "escrow": 0,
                    "available": 0,
                    "withdrawn": 0,
                    "disputed": 0,
                    "refunded": 0,
                    "inaccurate": false
                }
            },
            "min_withdraw": {
                "USD": 5.00,
                "ZWL": 2000,
            }             
        }

        v1.wallet.summary(false).then(summary => {
            this.setState({summary})
        }).catch(err => {
            return v1.wallet.summary(true).then(summary => {
                this.setState({ summary })
            })
        })

        this.updateSummary = (currency, summary) => {
            let curr = this.state.summary
            curr[currency] = summary
            this.setState({
                summary: curr
            })
        }

    }
    render() {
        return <>
            {this.props.user.verified ? (this.props.user.emailVerified ?
                <>
                    <div class="box">
                        <div class="box-header">
                            <h4 class="box-title">{this.state.currency} summary</h4>
                            <div class="box-controls pull-right">
                                <select class="form-control" onChange={(evt) => { this.setState({ currency: evt.target.value }) }}>
                                    <option value="ZWL">ZWL $</option>
                                    <option value="USD">USD $</option>
                                </select>
                            </div>
                        </div>

                        <div class="box-body">
                            {this.props.hideSummary ? <></> : <WalletSummary currency={this.state.currency} summary={this.state.summary[this.state.currency]} />}
                            {this.props.hideRecent ? <></> : <RecentWalletOperations currency={this.state.currency} />}
                            <hr />
                            {this.props.hideWithdraw ? <></> : <WalletRequestWithdrawal currency={this.state.currency} onWithdraw={this.updateSummary} min_withdraw={this.state.min_withdraw[this.state.currency]} summary={this.state.summary[this.state.currency]} />}
                            <hr />
                            {this.props.hidePayout ? <></> : <WalletUpdatePayoutDetails currency={this.state.currency} />}
                            <hr />
                            {this.props.showApiSettings ? <WalletAPISettings api={this.state.summary["api"]} /> : <></>}
                        </div>
                    </div>
                
                
                </> : <ActionCenter user={this.props.user} action='email' />
            ) : <Redirect to="/creator/wallet/mypayments" />}
        
        </>;
    }
}

export default Wallet;