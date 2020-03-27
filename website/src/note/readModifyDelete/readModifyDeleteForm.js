import React from 'react';
import Form from "react-bootstrap/Form";
import OverlayTrigger from "react-bootstrap/OverlayTrigger";
import Tooltip from "react-bootstrap/Tooltip";
import Button from "react-bootstrap/Button";
import Spinner from "react-bootstrap/Spinner";
import Modal from "react-bootstrap/Modal";

const ReadModifyDeleteForm = (props) =>
    <div>
        <Form onSubmit={props.onModify}>
            <Form.Group controlId="c">
                <Form.Control ref={props.contentRef} defaultValue={props.defaultContent} as="textarea"
                              maxLength={140} minLength={1} required rows={6}
                              disabled={!props.writeAccessChecked} onChange={props.contentChange}/>
                {props.writeAccessChecked && <Form.Text className="text-muted characters-hint" float="right">
                    {props.chars} characters left
                </Form.Text>}
            </Form.Group>
            <OverlayTrigger
                placement="right"
                delay={{show: 250, hide: 400}}
                overlay={
                    <Tooltip id="button-tooltip">
                        maximum 1 day
                    </Tooltip>}>
                <Form.Group controlId="t">
                    <Form.Control ref={props.ttlRef} type="number" min={1} max={1440}
                                  defaultValue={props.defaultTTL} disabled={!props.writeAccessChecked}/>
                </Form.Group>
            </OverlayTrigger>
            <Form.Group controlId="formBasicCheckbox">
                <Form.Check className="text-muted characters-hint" onChange={props.checkChange} checked={props.writeAccessChecked} type="checkbox"
                            label="I have write password"/>
            </Form.Group>
            {props.writeAccessChecked && <Form.Group controlId="passwordForm">
                <Form.Control type="password" ref={props.passwordRef} maxLength={50}
                              onFocus={props.onPasswordFocus} placeholder="write password"
                />
            </Form.Group>}
            {props.modifyError && <Form.Text className="passwordHint">
                incorrect password or note does not exist
            </Form.Text>}
            <Modal.Footer>
                {props.writeAccessChecked &&
                <>
                    <Button type="button"
                            variant='danger'
                            disabled={props.deleteSpinner || props.modifySpinner}
                            onClick={props.onDelete}>
                        {props.deleteSpinner && <Spinner
                            as="span"
                            animation="grow"
                            size="sm"
                            role="status"
                            aria-hidden="true"
                        />}
                        Delete
                    </Button>
                    <Button type="submit"
                            variant='primary'
                            disabled={props.deleteSpinner || props.modifySpinner}>
                        {props.modifySpinner && <Spinner
                            as="span"
                            animation="grow"
                            size="sm"
                            role="status"
                            aria-hidden="true"
                        />}
                        Modify
                    </Button>
                </>}
            </Modal.Footer>
        </Form>
    </div>;
export default ReadModifyDeleteForm;