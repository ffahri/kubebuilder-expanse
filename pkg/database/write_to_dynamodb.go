package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"simplekubebuilder/api/v1beta1"
)

const TABLE_NAME = "kubebuilder-expanse"

type DBConfig struct {
	DynamoDB *dynamodb.DynamoDB
}

func InitDynamoDB() *dynamodb.DynamoDB {
	mySession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")}))
	return dynamodb.New(mySession)

}

func (d *DBConfig) Init() {
	d.DynamoDB = InitDynamoDB()
}

func (d *DBConfig) Write(ships v1beta1.SpaceShips) error {
	putItem := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(ships.Spec.Name),
			},
			"class": {
				S: aws.String(ships.Spec.Class),
			},
			"owner": {
				S: aws.String(ships.Spec.Owner),
			},
			"status": {
				S: aws.String(string(ships.Status.Phase)),
			},
		},
		TableName: aws.String(TABLE_NAME),
	}
	_, err := d.DynamoDB.PutItem(putItem)
	if err != nil {
		return err
	}
	return nil
}
func (d *DBConfig) Update(ships v1beta1.SpaceShips) error {

	_, err := d.DynamoDB.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#s": aws.String("status"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {
				S: aws.String(string(ships.Status.Phase)),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(ships.Spec.Name),
			},
		},
		TableName:        aws.String(TABLE_NAME),
		UpdateExpression: aws.String("SET #s=:status"),
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *DBConfig) Delete(ships v1beta1.SpaceShips) error {
	putItem := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(ships.Spec.Name),
			},
			"class": {
				S: aws.String(ships.Spec.Class),
			},
			"owner": {
				S: aws.String(ships.Spec.Owner),
			},
			"status": {
				S: aws.String(string(ships.Status.Phase)),
			},
		},
		TableName: aws.String(TABLE_NAME),
	}
	_, err := d.DynamoDB.PutItem(putItem)
	if err != nil {
		return err
	}
	return nil
}

func (d *DBConfig) Get() {

}
