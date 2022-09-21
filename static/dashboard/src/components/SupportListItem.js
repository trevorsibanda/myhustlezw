import React, {Component} from "react"
import {Link} from "react-router-dom"
import ReactTimeAgo from "react-time-ago"
import money from "./payments/Amount"

class SupporterListItem extends Component {
    render() {
        let support_action = 'subscribed to your content'
        let support_card = <span className="btn btn-block btn-primary btn-sm btn-rounded" href="#">Subscription</span>

        switch (this.props.supporter.support_type){
            case 'support':
                support_action = 'bought you '+ this.props.supporter.items + ' '+ this.props.supporter.item_name
                if (this.props.supporter.items > 1){
                    support_action = support_action+'s'
                }
                support_card = <span className="btn btn-block btn-info btn-sm btn-rounded" href="#">Support</span>
                break;
            case 'subscribed':
                support_action = 'subscribed to  your account'
                support_card = <span className="btn btn-block btn-success btn-sm btn-rounded" href="#">Subscriber</span>
                break;
            case 'service_request':
                support_action = 'placed an order for "'+ this.props.supporter.item_name + '" '
                
                support_card = <span className="btn btn-block btn-info btn-sm btn-rounded" href="#">Order</span>
            break;
            case 'paid_content':
                support_action = 'paid to view "'+ this.props.supporter.item_name+'"'
                support_card = <span className="btn btn-block btn-info btn-sm btn-rounded" href="#">PayToView</span>
                break;
            default:
                support_action = 'Unknown support action'
                support_card = <span className="btn btn-block btn-danger btn-sm btn-rounded" href="#">Unknown</span>
        }
        return (
            <div className="col-md-6 col-12">
                <Link to={ "/creator/supporters/"+ this.props.supporter._id} className="media align-items-center bg-white mb-20">
                    <img className="avatar" src="/assets/img/placeholder.png" alt="..." />
                    <div className="media-body">
                        <p><strong>{this.props.supporter.display_name}</strong></p>
                        <p>{support_action}</p>
                        <p><ReactTimeAgo date={this.props.supporter.created_at} /></p>
                    </div>
                    <div>
                        <span className="btn btn-block btn-default btn-sm btn-rounded" href="/creator/wallet">{money.format(this.props.supporter.currency, this.props.supporter.amount)}</span>
                        {support_card}
                        
                    </div>
                </Link>
            </div>
        )
    }
}

export default SupporterListItem;
