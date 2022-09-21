import React, { Component } from "react"
import { Link, Switch, Route } from 'react-router-dom'

import './common.css'


class CreatorPublicSupporter extends Component {
    
    constructor(props) {
        super(props)

        this.actionOnClick = this.actionOnClick.bind(this)

    }

    actionOnClick(evt){
        if (this.props.supporter.creator_id === this.props.user._id) {
            return
        }
        evt.preventDefault()        
    }

    render() {

        let props = this.props
        var message
        var icon
        var who = 'you'
        console.log('here')
        console.log(props)
        if (props.supporter.creator_id !== props.user._id) {
            who = '@' + props.creator.username
        }
        switch (props.supporter.support_type) {
            case "subscribed": message = " subscribed to @" + props.creator.username + "'s account"; break;
            case "paid_content": message = " paid to view " + props.supporter.item_name; break;
            case "support": message = " bought " + who + " " + props.supporter.items + " " + props.supporter.item_name; icon = props.supporter.item_name; break;
            case "service_request": message = " placed an order for '" + props.supporter.item_name + "'"; break;
            default: message = "contributed to @" + props.creator.username + "'s account";
        }


        return (
            <Link class={"list-item " + (props.showMobileOnly ? 'd-none d-md-block' : '')} onClick={this.actionOnClick} to={'/creator/supporters/' + props.supporter._id} >
                <div><span>
                    <span class="w-48 avatar gd-primary"><img src="https://i.pinimg.com/564x/b6/4a/e2/b64ae23afaba459b30c3fcc0ad9e0d05.jpg" /></span>
                </span></div>
                <div class="flex">
                    <span class="item-author text-color" data-abc="true"><b>{props.supporter.display_name}</b></span>
                
                    {props.supporter.comment !== "" ?
                    
                        <div class="alert alert-info">
                            <div class="item-except text-muted text-sm h-1x">{message}</div>
                            <p>{props.supporter.comment.substr(0, 50)}{props.supporter.comment.length > 50 ? <a href="#">...</a> : <></>}</p>
                        </div> :
                    
                        <div class="item-except text-muted text-sm h-1x">{message}</div>}
                </div>
            </Link>
        )
    }
}

function CreatorRecentSupporters(props) {
    let grandMax = props.grandMax ? props.grandMax : props.supporters.length
    return (
<div class="">
    <h6>Recent supporters</h6>
    <hr class="hr-success" />
    <div class="list list-row block">
                {props.supporters.slice(0, grandMax).map((supporter, idx) => (
                   <CreatorPublicSupporter user={props.user} showMobileOnly={idx+1 > (props.maxShowMobile ? props.maxShowMobile : 200)} creator={props.creator} supporter={supporter} />
                ))}
            </div>
    {props.supporters.length > props.maxShowMobile ? 
            <div class="row justify-content-center" >
                <div class="col-lg-7 col-xs-12 col-md-7 order-md-1">
                    <div class="text-center">
                        <Link to={`/@${props.creator.username}/supporters`} class="btn btn-info btn-block"><i class="fa fa-heart" style={{color:'pink'}}></i> View more supporters</Link>
                    </div>
                </div>
            </div>  : <></>}
</div>
    )
}

export default CreatorRecentSupporters;