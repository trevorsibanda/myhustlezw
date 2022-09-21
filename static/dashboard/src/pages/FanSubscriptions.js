import React, { Component } from "react";
import v1 from "../api/v1";

import OptionsButtonGroup from "../components/OptionsButtongroup";
import Preloader from "../components/PreLoader";

import SubscriptionList from "../components/SubscriptionList";

class FanSubscriptions extends Component {

    constructor(props) {
        super(props)

        this.state = {
            subscriptions: [],
            filter: 'all',
            loading: true,
            viewOptions: [
            {
                value: 'all',
                component: <><i class="fa fa-eye"></i> All</> 
                },
                {
                value: 'active',
                component: <><i class="fa fa-magic"></i> Active Only</> 
                },
                {
                value: 'expired',
                component: <><i class="fa fa-history"></i> Expired</> 
                },
            ],
        }

        this.reloadFilter = (value) => {
            this.setState({ filter: value, loading: true })
            v1.subscriptions.listAll(value, false).then(subscriptions => {
                this.setState({ subscriptions, loading: false })
            }).catch(err => {
                alert('Failed to apply filter with error ' + (err.error ? err.error : err))
                this.setState({loading: false})
            })
        }

        v1.subscriptions.listAll(this.state.filter, false).then(subscriptions => {
            this.setState({ subscriptions, loading: false })
        }).catch(_ => {
            v1.subscriptions.listAll(this.state.filter, true).then(subscriptions => {
                this.setState({ subscriptions })
            })
        })
    }

    render() {

        return (
            <>
                <div class="box">
                    <div class="box-header">
                        <div class="row">
                            <div class="col-md-8">
                                <h4 class="box-title">Your Subscriptions</h4>
                            </div>
                            <div class="col-md-4">
                                <div class="box-controls pull-right">
                                    <OptionsButtonGroup items={this.state.viewOptions} item={this.state.filter} onChange={this.reloadFilter} />
                                    
                                &nbsp;
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                {this.state.loading ? <Preloader /> : <SubscriptionList subscriptions={this.state.subscriptions} fans={true} />}
                
            </>
        )
    }

}

export default FanSubscriptions;