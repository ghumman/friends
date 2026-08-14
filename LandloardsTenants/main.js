class ApplicationEvaluatorResponse {
  constructor(input) {
    this.acceptableApplications = [];
    this.rejectedApplications = [];
    
    if (input) {
      const { acceptableApplications, rejectedApplications } = input;
      this.acceptableApplications = acceptableApplications;
      this.rejectedApplications = rejectedApplications;
    }
  }
}

class RejectedApplication {
  constructor(applicationId, rejectedReasons) {
    this.applicationId = applicationId;
    this.rejectedReasons = rejectedReasons;
  }
}

class TestInput {
  constructor(applications, incomeSources, landlordRequirements) {
    this.applications = applications;
    this.incomeSources = incomeSources;
    this.landlordRequirements = landlordRequirements;
  }
}

class Application {
  constructor(applicationId, firstName, lastName, email, phoneNumber, userId, creditScore) {
    this.applicationId = applicationId;
    this.firstName = firstName;
    this.lastName = lastName;
    this.email = email;
    this.phoneNumber = phoneNumber;
    this.userId = userId;
    this.creditScore = creditScore;
  }
}

class IncomeSource {
  constructor(title, employer, grossIncome, netIncome, startDate, endDate, userId, periodType) {
    this.title = title;
    this.employer = employer;
    this.grossIncome = grossIncome;
    this.netIncome = netIncome;
    this.startDate = new Date(startDate);
    this.endDate = endDate ? new Date(endDate) : null;
    this.userId = userId;
    this.periodType = periodType;
  }
}

class LandlordRequirements {
  constructor(monthlyRent, minimumIncomeToRentRatio, minimumCreditScore) {
    this.monthlyRent = monthlyRent;
    this.minimumIncomeToRentRatio = minimumIncomeToRentRatio;
    this.minimumCreditScore = minimumCreditScore;
  }
}

const IncomePeriodType = Object.freeze({
  MONTHLY: 'Monthly',
  YEARLY: 'Yearly'
});

const StringToIncomePeriodType = Object.freeze({
  'Monthly': IncomePeriodType.MONTHLY,
  'Yearly': IncomePeriodType.YEARLY,
})

const RejectionReason = {
  INCOME_RENT_RATIO_NOT_MET: "Applicant does not meet the required income to rent ratio.",
  INVALID_CREDIT_SCORE_RANGE: "Credit score did not fall within expected range [300-850]",
  NULL_OR_EMPTY_EMAIL: "Email was null or empty.",
  MINIMUM_CREDIT_SCORE_NOT_MET: "Credit Score did not meet minimum requirement"
};

function evaluateApplications(testInput) {
  // your code here
  
  let rejectedApp = []
  let acceptedApp = []
  
  for (const a of testInput.applications) {
    
    let rejectionReasons = []
    if (a.email ==  null || a.email == "") {
      rejectionReasons.push(RejectedApplication.NULL_OR_EMPTY_EMAIL)
      
    } 
    if (a.creditScore < 300 || a.creditScore >850 ) {
      rejectionReasons.push(RejectedApplication.INVALID_CREDIT_SCORE_RANGE)
    }
    
    if (a.creditScore < testInput.landlordRequirements.minimumCreditScore) {
      rejectionReasons.push(RejectedApplication.MINIMUM_CREDIT_SCORE_NOT_MET)
    }
    
    let totalIncome = 0; 
    
    for (const i of testInput.incomeSources) {
      if (i.userId == a.userId) {
        totalIncome = totalIncome +  i.netIncome;
      }
    }
    
    
    const ratio = totalIncome / testInput.landlordRequirements.monthlyRent;
    
    if (ratio < testInput.landlordRequirements.minimumIncomeToRentRatio) {
      rejectionReasons.push(RejectedApplication.INCOME_RENT_RATIO_NOT_MET)
    }
    
    if (rejectionReasons.length != 0) {
      rejectedApp.push(new RejectedApplication(
        a.applicationId,
        rejectionReasons
      ))
    } else {
      acceptedApp.push(a.applicationId)
    }
  }
  
  
  return new ApplicationEvaluatorResponse({acceptableApplications: acceptedApp, rejectedApplications: rejectedApp});
}

function constructTestInput(jsonObj) {
  let inputApplications = []
  let inputIncomeSources = []
  
  for (const a of jsonObj.applications) {
    inputApplications.push(new Application(
      a.applicationId,
      a.firstName,
      a.lastName,
      a.email,
      a.phoneNumber,
      a.userId,
      a.creditScore,
    ))
  }
  
  for (const i of jsonObj.incomeSources) {
    inputIncomeSources.push(new IncomeSource(
      i.title,
      i.employer,
      i.grossIncome,
      i.netIncome,
      i.startDate,
      i.endDate,
      i.userId,
      StringToIncomePeriodType[i.periodType],
    ))
  }
  
  let inputLandlordRequirements = null;
  if (jsonObj.landlordRequirements) {
    req = jsonObj.landlordRequirements
    inputLandlordRequirements = new LandlordRequirements(
      req.monthlyRent,
      req.minimumIncomeToRentRatio,
      req.minimumCreditScore,
    )
  }
  
  return new TestInput(inputApplications, inputIncomeSources, inputLandlordRequirements)
}

module.exports = { ApplicationEvaluatorResponse, RejectedApplication, RejectionReason, evaluateApplications, constructTestInput };