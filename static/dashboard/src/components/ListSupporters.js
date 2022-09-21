import React, {Component} from "react"
import SupporterListItem from "./SupportListItem"
import { Link } from "react-router-dom"
import v1 from "../api/v1"



class SupportersList extends React.Component {
    constructor(props) {
        super(props)


        this.state = {
            supporters: [],
            loading: true,
            skipCounter: 0,
        }
        let promise = null


        switch (this.props.type) {
            case 'recent':
                promise = v1.supporters.recent(this.props.count ? this.props.count : 50, 0, false).catch(_ => {
                    return v1.supporters.recent(this.props.count ? this.props.count : 10, 0, true)
                })
            break;
            case 'recent_campaign':
                promise = v1.supporters.recent_campaign(this.props.campaign_id, this.props.count ? this.props.count : 50, 0 , false).catch(_ => {
                    return v1.supporters.recent(this.props.count ? this.props.count : 10, 0, true)
                })
                break;
            case 'all':
                promise = v1.supporters.all(false).catch(_ => v1.supporters.all(true))
            break;
        }


        promise.then(resp => {
            if(resp.supporters){
                this.setState({ supporters: resp.supporters, loading: false })
            }
            
        }).catch(err => {
            this.setState({
                err: 'Failed to load supporters'
            })
        })
    }

    render() {

        let supportersList = (
        <div className={this.state.supporters.length === 0 ? "d-none " : "col-md-12 col-12"}>
            <div className="media flex-column text-center p-40 bg-white mb-20">
                <span className="avatar avatar-xxl bg-white opacity-60 mx-auto">
                    <i className="align-sub fa fa-frown bg-white font-size-40"></i>
                </span>
                <div className="mt-20">
                    <h6 className="text-uppercase fw-500">Seems you have no supporters yet!</h6>
                    <small>Share your link to your fans, clients and supporters. A good place is social media.</small>
                </div>
            </div>
        </div>)
        if (this.state.supporters.length > 0) {
            supportersList = this.state.supporters.map((item, index) => {
                    return <SupporterListItem index={index} supporter={item} />
            })
        }


        return (
            <>
                <h5>{this.props.title}</h5>
                <div className="row" >
                    {supportersList}
                </div>
                <div className="row justify-content-center">
                    <div className="col-md-3">
                        <Link to="/creator/supporters" className="btn btn-default btn-block"><i className="fa fa-heart"></i> See all {this.props.supporterName}s</Link>
                    </div>
                </div>
            </>
        )
    }
}

export default SupportersList;