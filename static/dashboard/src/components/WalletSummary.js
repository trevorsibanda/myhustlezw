import React, {Component} from "react"
import money from "./payments/Amount"

class WalletSummary extends Component {
    render() {
        return (
            <>
                <div class="row">

                    <div class="col-lg-3 col-md-6">
                        <div class="box bg-gradient-deepocean box-inverse">
                            <div class="box-body">
                                <h3>{money.format(this.props.currency, this.props.summary.withdrawn)}</h3>
                                <span class="pull-right">Total earned</span>
                            </div>
                        </div>
                    </div>
                    
                    <div class="col-lg-3 col-md-6">
                        <div class="box bg-success box-inverse">
                            <div class="box-body ">
                                <h3>{money.format(this.props.currency, this.props.summary.available)}</h3>
                                <span class="pull-right">Ready to withdraw</span>
                            </div>
                        </div>
                    </div>
                    { this.props.summary.pending_withdrawal <= 0 ? <></> :
                    <div class="col-lg-3 col-md-6">
                        <div class="box bg-danger box-inverse">
                            <div class="box-body ">
                                <h3>{money.format(this.props.currency, this.props.summary.pending_withdrawal)}</h3>
                                <span class="pull-right">Processing withdrawal</span>
                            </div>
                        </div>
                    </div> 
                    }
                    <div class="col-lg-3 col-md-6">
                        <div class="box bg-gradient-purple box-inverse">
                            <div class="box-body">
                                <h3>{money.format(this.props.currency, this.props.summary.escrow)} </h3>
                                <span class="pull-right">In escrow</span>
                            </div>
                        </div>
                    </div>
                    <div class="col-lg-3 col-md-6">
                        <div class="box bg-gradient-botani box-inverse">
                            <div class="box-body">
                                <h3>{money.format(this.props.currency, this.props.summary.disputed)}</h3>
                                <span class="pull-right">Disputed</span>
                            </div>
                        </div>
                    </div>
                </div>
                <hr />
            </>
        )
    }
}


export default WalletSummary;