package usecase

import (
	"ProjectMessenger/domain"
	"ProjectMessenger/microservices/contacts_service/proto"
	"context"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ContactStore interface {
	GetContacts(ctx context.Context, userID uint) []domain.Person
	AddContact(ctx context.Context, userID1, userID2 uint) (ok bool)
}

type ContactsManager struct {
	chats.UnimplementedContactsServer
	storage ContactStore
}

func NewContactsManager(storage ContactStore) *ContactsManager {
	return &ContactsManager{
		storage: storage,
	}
}

func convertToNormalUser(person *chats.Person) domain.Person {
	return domain.Person{
		ID:           uint(person.GetID()),
		Username:     person.GetUsername(),
		Email:        person.GetEmail(),
		Name:         person.GetName(),
		Surname:      person.GetSurname(),
		About:        person.GetAbout(),
		Password:     person.GetPassword(),
		CreateDate:   person.CreateTime.AsTime(),
		LastSeenDate: person.LastSeenDate.AsTime(),
		AvatarPath:   person.GetAvatarPath(),
		PasswordSalt: person.GetPasswordSalt(),
	}
}

func convertToGRPCUser(person domain.Person) *chats.Person {
	return &chats.Person{
		ID:           uint64(person.ID),
		Username:     person.Username,
		Email:        person.Email,
		Name:         person.Name,
		Surname:      person.Surname,
		About:        person.About,
		Password:     person.Password,
		CreateTime:   timestamppb.New(person.CreateDate),
		LastSeenDate: timestamppb.New(person.LastSeenDate),
		AvatarPath:   person.AvatarPath,
		PasswordSalt: person.PasswordSalt,
	}
}

func (cm *ContactsManager) GetContacts(ctx context.Context, in *chats.UserIDContacts) (*chats.PersonArray, error) {
	userID := uint(in.GetUserID())
	contactsFromStorage := cm.storage.GetContacts(ctx, userID)
	resp := &chats.PersonArray{Persons: make([]*chats.Person, 0)}
	for i := range contactsFromStorage {
		resp.Persons = append(resp.Persons, convertToGRPCUser(contactsFromStorage[i]))
	}
	return resp, nil
}

func (cm *ContactsManager) AddContactByUsername(ctx context.Context, in *chats.AddByUsernameReq) (*chats.EmptyContacts, error) {
	usernameToAdd := in.GetUsernameToAdd()
	userAddingID := uint(in.GetUserAddingID())
	userToAddID := uint(in.GetUserToAddID())
	contactsFromStorage := cm.storage.GetContacts(ctx, userAddingID)
	for i := range contactsFromStorage {
		if contactsFromStorage[i].Username == usernameToAdd {
			return &chats.EmptyContacts{}, status.Error(400, "Такой контакт уже существует")
		}
	}

	ok := cm.storage.AddContact(ctx, userAddingID, userToAddID)
	if !ok {
		return &chats.EmptyContacts{}, status.Error(500, "")
	}
	return &chats.EmptyContacts{}, nil
}

func (cm *ContactsManager) AddToAllContacts(ctx context.Context, in *chats.AddToAllReq) (*chats.BoolResponseContacts, error) {
	userAddingID := uint(in.GetUserAddingID())
	userIDsRPC := in.Users.Users
	userIDs := make([]uint, 0)
	for i := range userIDsRPC {
		userIDs = append(userIDs, uint(userIDsRPC[i].GetUserID()))
	}
	for i := range userIDs {
		if userIDs[i] == userAddingID {
			continue
		}
		if !cm.storage.AddContact(ctx, userAddingID, userIDs[i]) {
			return &chats.BoolResponseContacts{Ok: false}, nil
		}
	}
	return &chats.BoolResponseContacts{Ok: true}, nil
}
