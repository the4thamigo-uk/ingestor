package contact

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TODO: tests not complete

func TestParser(t *testing.T) {
	data := bytes.NewBuffer([]byte(testData))
	r, err := NewReader(data)
	require.Nil(t, err)
	cs, err := r.ReadAll()
	require.Nil(t, err)
	assert.Len(t, cs, 10)
	assert.Equal(t, "1", cs[0].ID)
	assert.Equal(t, "Kirk", cs[0].Name)
	assert.Equal(t, "ornare@sedtortor.net", cs[0].Email)
	assert.Equal(t, "+441389037420", cs[0].Mobile)
}

func TestParser_ReversedFields(t *testing.T) {
	data1 := bytes.NewBuffer([]byte(testData))
	r1, err := NewReader(data1)
	require.Nil(t, err)
	cs1, err := r1.ReadAll()
	require.Nil(t, err)

	data2 := bytes.NewBuffer([]byte(testData2))
	r2, err := NewReader(data2)
	require.Nil(t, err)
	cs2, err := r2.ReadAll()
	require.Nil(t, err)

	assert.Equal(t, cs1, cs2)
}

func TestParser_NoDataShouldFail(t *testing.T) {
	data := bytes.NewBuffer(nil)
	r, err := NewReader(data)
	assert.NotNil(t, err)
	assert.Nil(t, r)
}

func TestParser_OnlyHeaderReturnsNoData(t *testing.T) {
	data := bytes.NewBuffer([]byte(`id,name,email,mobile_number`))
	r, err := NewReader(data)
	require.Nil(t, err)
	cs, err := r.ReadAll()
	require.Nil(t, err)
	assert.Len(t, cs, 0)
}

var testData = `
id,name,email,mobile_number
1,Kirk,ornare@sedtortor.net,(013890) 37420
2,Cain,volutpat@semmollisdui.com,(016977) 2245
3,Geoffrey,vitae@consectetuermaurisid.co.uk,0800 1111
4,Walter,odio.a.purus@sit.edu,(0161) 328 6656
5,Armand,Cras.vulputate@metusvitae.co.uk,0836 796 0064
6,Travis,malesuada.id.erat@faucibusid.com,07186 118681
7,Christian,amet@etipsum.edu,(0116) 453 8054
8,Warren,amet.luctus.vulputate@hendreritnequeIn.co.uk,0891 091 2561
9,Ian,Fusce@elitpharetra.org,0845 46 41
10,Cairo,auctor.vitae@tristiqueac.org,07624 856071
`
var testData2 = `
mobile_number,email,name,id
(013890) 37420,ornare@sedtortor.net,Kirk,1
(016977) 2245,volutpat@semmollisdui.com,Cain,2
0800 1111,vitae@consectetuermaurisid.co.uk,Geoffrey,3
(0161) 328 6656,odio.a.purus@sit.edu,Walter,4
0836 796 0064,Cras.vulputate@metusvitae.co.uk,Armand,5
07186 118681,malesuada.id.erat@faucibusid.com,Travis,6
(0116) 453 8054,amet@etipsum.edu,Christian,7
0891 091 2561,amet.luctus.vulputate@hendreritnequeIn.co.uk,Warren,8
0845 46 41,Fusce@elitpharetra.org,Ian,9
07624 856071,auctor.vitae@tristiqueac.org,Cairo,10
`
