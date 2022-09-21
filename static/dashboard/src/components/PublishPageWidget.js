
import React, {Component} from "react"
import v1 from "../api/v1"

class PublishPageWidget extends Component {

    publishSite() {
        v1.user.publishPage(!this.props.user.page.published).then(resp => {
            if (resp.status === 'ok') {
                alert('Operation completed successfully. Will now reload page.')
                window.reload(true)
            } else {
                //notify of error
                alert('Encountered error : ' + resp.error)
            }
        }).catch(err => {
            alert('Failed to publish site with error: ' + err)
        })
    }
    render() {
        return (
            <div className="box">
                <div className="box-header">
                    <h4 className="box-title">Publish your page</h4>
                    <div className="box-controls pull-right">
                        {
                            this.props.user.page.published ?
                                <button className="btn btn-xs btn-success btn-disabled">Page is Live!</button> :
                                <button className="btn btn-xs btn-danger btn-disabled" >Page is not yet published!</button>
                        }

                    </div>
                </div>
                <div className="box-body">
                    <p>{
                        this.props.user.page.published ?
                            <>Your page is currently published and viewable at
                        <a href={'/@' + this.props.user.username} target="_blank">https://myhustle.co.zw/@{this.props.user.username}</a>
                            </> : <>Your page and all your content is currently not published. </>
                    }</p>

                    {
                        this.props.user.page.published ?
                            <a className="btn btn-outline btn-danger btn-block text-white" onClick={this.publishSite}><i className="fa fa-eye"></i> Make site private.</a> :
                            <a className="btn btn-info btn-block text-white" disabled={!this.props.user.phoneVerified} onClick={this.publishSite}><i className="fa fa-check"></i> Publish your page.</a>
                    }

                    {!this.props.user.phoneVerified ?
                        <strong>You must verify your phone number or email address first</strong> :
                        <></>
                    }

                </div>
            </div>
        )
    }
}

export default PublishPageWidget;